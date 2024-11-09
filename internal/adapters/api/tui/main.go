package tui

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"syscall"
	"time"

	gclient "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
	"github.com/k0st1a/gophkeeper/internal/adapters/api/tui/storage"
	"github.com/k0st1a/gophkeeper/internal/pkg/client/model"
	"github.com/k0st1a/gophkeeper/internal/pkg/job"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	// Имена страниц.
	pageNameWelcome  = "welcome"
	pageNameRegister = "register"
	pageNameLogin    = "login"
	pageNameItems    = "items"
	pageNameNotify   = "notify"

	pageNameUpdatePassword = "update password"
	pageNameAddPassword    = "add password"

	pageNameUpdateCard = "update card"
	pageNameAddCard    = "add card"

	pageNameUpdateNote = "update note"
	pageNameAddNote    = "add note"

	pageNameUpdateFile = "update file"
	pageNameAddFile    = "add file"

	// Имена кнопок.
	buttonNameCancel = "Cancel"
	buttonNameOk     = "Ok"
	buttonNameUpdate = "Update"
	buttonNameDelete = "Delete"

	// Имена надписей.
	labelName                  = "Name"
	labelDescription           = "Description"
	labelResource              = "Resource"
	labelUserName              = "User name"
	labelPassword              = "Password"
	labelAdditionalInformation = "Additional information"
	labelCardNumber            = "Card Number"
	labelCardExpires           = "Card expires"
	labelCardHolder            = "Card holder"
	labelNote                  = "Note"
	labelAdd                   = "Add"

	defaultFieldWidth  = 30
	defaultFieldHeight = 5
	defaultMaxLength   = 255
)

const (
	columnName = iota
	columnType
	columnUpdateTime
)

type client struct {
	grpc    gclient.UserAuthentication
	storage storage.ItemStorage
	sync    job.StartStopper
	cancel  func()
	app     *tview.Application
	pages   *tview.Pages
}

func New(c gclient.UserAuthentication, s storage.ItemStorage, j job.StartStopper, cn func()) *client {
	app := tview.NewApplication()
	pages := tview.NewPages()

	app.SetRoot(pages, true).EnableMouse(true)

	return &client{
		grpc:    c,
		storage: s,
		sync:    j,
		cancel:  cn,
		app:     app,
		pages:   pages,
	}
}

func (c *client) Run(ctx context.Context) error {
	log.Ctx(ctx).Printf("Run tui")
	c.WelcomePage(ctx)
	err := c.app.Run()
	if err != nil {
		return fmt.Errorf("error of run tui client:%w", err)
	}

	return nil
}

func (c *client) Stop(ctx context.Context) {
	log.Ctx(ctx).Printf("Stop tui")
	c.app.Stop()
	c.StopSync(ctx)
	c.cancel()
}

func (c *client) WelcomePage(ctx context.Context) {
	log.Printf("Invoked Welcome Page")

	welcomeList := tview.NewList().
		ShowSecondaryText(false).
		AddItem("Login", "", '1', func() {
			c.pages.RemovePage(pageNameWelcome)
			c.LoginPage(ctx)
		}).
		AddItem("Register", "", '2', func() {
			c.pages.RemovePage(pageNameWelcome)
			c.RegisterPage(ctx)
		}).
		AddItem("Quit", "", 'q', func() {
			c.Stop(ctx)
		})

	welcomeList.
		SetTitle("Welcome").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	welcomeFlexBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(welcomeList, 0, 1, true)

	c.pages.AddPage(pageNameWelcome, welcomeFlexBox, true, true)
}

func (c *client) NotifyPage(text string) {
	log.Printf("Invoked Notify Page, text:%v", text)
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{buttonNameCancel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == buttonNameCancel {
				c.pages.RemovePage(pageNameNotify)
			}
		})

	c.pages.AddPage(pageNameNotify, modal, true, true)
}

func (c *client) NotifyAndSwitch2Page(text string, page func()) {
	log.Printf("Invoked Notify Page, text:%v", text)
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{buttonNameOk}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == buttonNameOk {
				c.pages.RemovePage(pageNameNotify)
				page()
			}
		})

	c.pages.AddPage(pageNameNotify, modal, true, true)
}

func (c *client) RegisterPage(ctx context.Context) {
	log.Printf("Invoked Register Page")

	var (
		email           string
		password        string
		confirmPassword string
	)
	registerForm := tview.NewForm().
		AddInputField("Email", "", defaultFieldWidth, nil, func(text string) {
			email = text
		}).
		AddPasswordField("Password", "", defaultFieldWidth, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Confirm password", "", defaultFieldWidth, '*', func(text string) {
			confirmPassword = text
		}).
		AddButton("Register", func() {
			if email == "" {
				c.NotifyPage("Email is empty")
				return
			}

			if password == "" {
				c.NotifyPage("Password is empty")
				return
			}

			if confirmPassword == "" {
				c.NotifyPage("Confirm password is empty")
				return
			}

			if password != confirmPassword {
				c.NotifyPage("Passwords not equals")
				return
			}

			err := c.grpc.RegisterUser(ctx, email, password)
			if err != nil {
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameRegister)

			c.NotifyAndSwitch2Page("Success register", func() {
				c.WelcomePage(ctx)
			})
		}).
		AddButton("Cancel", func() {
			c.pages.RemovePage(pageNameRegister)
			c.WelcomePage(ctx)
		})

	registerForm.
		SetTitle("Enter register data").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	registerFlexBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(registerForm, 0, 1, true)

	c.pages.AddPage(pageNameRegister, registerFlexBox, true, true)
}

func (c *client) LoginPage(ctx context.Context) {
	log.Printf("Invoked Login Page")

	var (
		email    string
		password string
	)
	loginForm := tview.NewForm().
		AddInputField("Email", "", defaultFieldWidth, nil, func(text string) {
			email = text
		}).
		AddPasswordField("Password", "", defaultFieldWidth, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			if email == "" {
				c.NotifyPage("Email is empty")
				return
			}

			if password == "" {
				c.NotifyPage("Password is empty")
				return
			}

			err := c.grpc.LoginUser(ctx, email, password)
			if err != nil {
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameLogin)
			log.Printf("Success login fast")

			c.NotifyAndSwitch2Page("Success login", func() {
				c.pages.RemovePage(pageNameLogin)
				c.StartSync(ctx)
				c.ItemsPage(ctx)
			})
			log.Printf("Success login fast")
		}).
		AddButton("Cancel", func() {
			c.pages.RemovePage(pageNameLogin)
			c.WelcomePage(ctx)
		})

	loginForm.
		SetTitle("Enter login data").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	loginFlexBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginForm, 0, 1, true)

	c.pages.AddPage(pageNameLogin, loginFlexBox, true, true)
}

func (c *client) ItemsPage(ctx context.Context) {
	log.Printf("Invoked Items Page")

	l, err := c.storage.ListItems(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error of list local items")
	}

	sort.Slice(l, func(i, j int) bool {
		return l[i].CreateTime.Before(l[j].CreateTime)
	})

	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetSeparator(' ').
		SetCell(0, columnName, tview.NewTableCell("Name").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnType, tview.NewTableCell("Type").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnUpdateTime, tview.NewTableCell("Update time").SetSelectable(false).SetTextColor(tcell.ColorYellow))

	table.
		SetBorder(true).
		SetTitle("Table")

	for i, item := range l {
		name, err := item.GetName()
		if err != nil {
			log.Error().Err(err).Msg("error of get item name")
			continue
		}

		itype, err := item.GetType()
		if err != nil {
			log.Error().Err(err).Msg("error of get item name")
			continue
		}

		row := i + 1

		table.
			SetCell(row, columnName, tview.NewTableCell(name).SetTextColor(tcell.ColorWhite).SetReference(item)).
			SetCell(row, columnType, tview.NewTableCell(itype).SetTextColor(tcell.ColorWhite)).
			SetCell(row, columnUpdateTime, newTableCellTime(item.UpdateTime).SetSelectable(false))
	}

	table.SetSelectedFunc(func(row, column int) {
		item, ok := table.GetCell(row, columnName).GetReference().(storage.Item)
		if !ok {
			log.Error().Msgf("error of get item by reference while selected, row:%v, column:%v", row, column)
			return
		}

		c.UpdateItemPage(ctx, &item)
	})

	table.
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	buttons := tview.NewForm().
		AddButton("Add password", func() {
			c.AddPasswordPage(ctx)
		}).
		AddButton("Add card", func() {
			c.AddCardPage(ctx)
		}).
		AddButton("Add note", func() {
			c.AddNotePage(ctx)
		}).
		AddButton("Add file", func() {
			c.AddFilePage(ctx)
		}).
		AddButton("Delete", func() {
			row, _ := table.GetSelection()
			item, ok := table.GetCell(row, columnName).GetReference().(storage.Item)
			if !ok {
				log.Error().Msgf("error of get item by reference while delete item, row:%v", row)
				return
			}

			name := table.GetCell(row, columnName).Text
			itype := table.GetCell(row, columnName).Text

			c.DeleteItemPage(ctx, &item, name, itype)
		}).
		AddButton("Refresh", func() {
			c.ItemsPage(ctx)
		}).
		AddButton("Logout", func() {
			c.StopSync(ctx)
			c.grpc.Logout(ctx)
			c.storage.Clear(ctx)
			c.pages.RemovePage(pageNameItems)
			c.WelcomePage(ctx)
		})

	buttons.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//nolint // no need here
		switch event.Key() {
		case tcell.KeyTab:
			c.app.SetFocus(table)
		case tcell.KeyRight:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyLeft:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})

	buttons.
		SetButtonsAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 0, 0)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(buttons, 1, 1, false).
		AddItem(table, 0, 1, true)

	flex.
		SetTitle("Items page").
		SetBorder(true)

	table.SetDoneFunc(func(key tcell.Key) {
		c.app.SetFocus(buttons)
	})

	c.pages.AddPage(pageNameItems, flex, true, true)
}

func (c *client) UpdatePasswordPage(ctx context.Context, i *storage.Item, p *storage.Password) {
	log.Printf("Invoked Updated password Page, item(%v)", i.ID)

	form := tview.NewForm().
		AddInputField(labelResource, p.Resource, defaultFieldWidth, nil, func(text string) {
			p.Resource = text
		}).
		AddInputField(labelUserName, p.UserName, defaultFieldWidth, nil, func(text string) {
			p.UserName = text
		}).
		AddInputField(labelPassword, p.Password, defaultFieldWidth, nil, func(text string) {
			p.Password = text
		}).
		AddTextArea(labelDescription, i.Meta.Get(model.MetaKeyDescription), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, i.Meta.Get(model.MetaKeyAdditionalInformation), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(buttonNameUpdate, func() {
			err := c.storage.UpdateItem(ctx, i)
			if err != nil {
				log.Error().Err(err).Msg("Item update error while update password")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameUpdatePassword)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameUpdatePassword)
		})

	form.
		SetTitle("Update password").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameUpdatePassword, flex, true, true)
}

func (c *client) AddPasswordPage(ctx context.Context) {
	log.Printf("Invoked Add password Page")

	p := &storage.Password{}
	m := storage.Meta{}

	//nolint:dupl // not need check for tview.NewForm
	form := tview.NewForm().
		AddInputField(labelResource, p.UserName, defaultFieldWidth, nil, func(text string) {
			p.Resource = text
		}).
		AddInputField(labelUserName, p.UserName, defaultFieldWidth, nil, func(text string) {
			p.UserName = text
		}).
		AddInputField(labelPassword, p.Password, defaultFieldWidth, nil, func(text string) {
			p.Password = text
		}).
		AddTextArea(labelDescription, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(buttonNameOk, func() {
			_, err := c.storage.CreateItem(ctx, p, m)
			if err != nil {
				log.Error().Err(err).Msg("Item add error while add password")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameAddPassword)
			c.ItemsPage(ctx)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameAddPassword)
		})

	form.
		SetTitle("Add password").
		SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameAddPassword, flex, true, true)
}

func (c *client) UpdateCardPage(ctx context.Context, i *storage.Item, cd *storage.Card) {
	log.Printf("Invoked Update card Page, item(%v)", i.ID)

	form := tview.NewForm().
		AddInputField(labelCardNumber, cd.Number, defaultFieldWidth, nil, func(text string) {
			cd.Number = text
		}).
		AddInputField(labelCardExpires, cd.Expires, defaultFieldWidth, nil, func(text string) {
			cd.Expires = text
		}).
		AddInputField(labelCardHolder, cd.Holder, defaultFieldWidth, nil, func(text string) {
			cd.Holder = text
		}).
		AddTextArea(labelDescription, i.Meta.Get(model.MetaKeyDescription), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, i.Meta.Get(model.MetaKeyAdditionalInformation), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(buttonNameUpdate, func() {
			err := c.storage.UpdateItem(ctx, i)
			if err != nil {
				log.Error().Err(err).Msg("Item update error while update card")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameUpdateCard)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameUpdateCard)
		})

	form.
		SetTitle("Update card").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameUpdateCard, flex, true, true)
}

func (c *client) AddCardPage(ctx context.Context) {
	log.Printf("Invoked Add card Page")

	cd := &storage.Card{}
	m := storage.Meta{}

	//nolint:dupl // not need check for tview.NewForm
	form := tview.NewForm().
		AddInputField(labelCardNumber, cd.Number, defaultFieldWidth, nil, func(text string) {
			cd.Number = text
		}).
		AddInputField(labelCardExpires, cd.Expires, defaultFieldWidth, nil, func(text string) {
			cd.Expires = text
		}).
		AddInputField(labelCardHolder, cd.Holder, defaultFieldWidth, nil, func(text string) {
			cd.Holder = text
		}).
		AddTextArea(labelDescription, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(buttonNameOk, func() {
			_, err := c.storage.CreateItem(ctx, cd, m)
			if err != nil {
				log.Error().Err(err).Msg("Item add error while add card")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameAddCard)
			c.ItemsPage(ctx)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameAddCard)
		})

	form.
		SetTitle("Add card").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameAddCard, flex, true, true)
}

func (c *client) UpdateNotePage(ctx context.Context, i *storage.Item, n *storage.Note) {
	log.Printf("Invoked Update note Page, item(%v)", i.ID)

	form := tview.NewForm().
		AddInputField(labelName, n.Name, defaultFieldWidth, nil, func(text string) {
			n.Name = text
		}).
		AddTextArea(labelNote, n.Body, defaultFieldWidth, 0, 0, func(text string) {
			n.Body = text
		}).
		AddTextArea(labelDescription, i.Meta.Get(model.MetaKeyDescription), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, i.Meta.Get(model.MetaKeyAdditionalInformation), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(buttonNameUpdate, func() {
			err := c.storage.UpdateItem(ctx, i)
			if err != nil {
				log.Error().Err(err).Msg("Item update error while update note")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameUpdateNote)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameUpdateNote)
		})

	form.
		SetTitle("Update note").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameUpdateNote, flex, true, true)
}

func (c *client) AddNotePage(ctx context.Context) {
	log.Printf("Invoked Add note Page")

	n := &storage.Note{}
	m := storage.Meta{}

	form := tview.NewForm().
		AddInputField(labelName, n.Name, defaultFieldWidth, nil, func(text string) {
			n.Name = text
		}).
		AddTextArea(labelNote, n.Body, defaultFieldWidth, 0, 0, func(text string) {
			n.Body = text
		}).
		AddTextArea(labelDescription, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(buttonNameOk, func() {
			_, err := c.storage.CreateItem(ctx, n, m)
			if err != nil {
				log.Error().Err(err).Msg("Item add error while add note")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameAddNote)
			c.ItemsPage(ctx)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameAddNote)
		})

	form.
		SetTitle("Add note").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameAddNote, flex, true, true)
}

func (c *client) UpdateFilePage(ctx context.Context, i *storage.Item, f *storage.File) {
	log.Printf("Invoked update file Page, item(%v)", i.ID)
	var path string

	form := tview.NewForm().
		AddInputField(labelName, f.Name, defaultFieldWidth, nil, func(text string) {
			f.Name = text
		}).
		AddInputField("Path to download", path, defaultFieldWidth, nil, func(text string) {
			path = text
		}).
		AddTextArea(labelDescription, i.Meta.Get(model.MetaKeyDescription), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, i.Meta.Get(model.MetaKeyAdditionalInformation), defaultFieldWidth,
			defaultFieldHeight, defaultMaxLength, func(text string) {
				i.Meta.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton("Download", func() {
			if err := os.WriteFile(path, f.Body, syscall.S_IRUSR|syscall.S_IWUSR); err != nil {
				c.NotifyPage(err.Error())
				return
			}
			c.pages.RemovePage(pageNameUpdateFile)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameUpdateFile)
		})

	form.
		SetTitle("Update file").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameUpdateFile, flex, true, true)
}

func (c *client) AddFilePage(ctx context.Context) {
	log.Printf("Invoked Add file Page")

	f := &storage.File{}
	m := storage.Meta{}

	var path string

	form := tview.NewForm().
		AddInputField("Path", path, defaultFieldWidth, nil, func(text string) {
			path = text
		}).
		AddTextArea(labelDescription, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyDescription, text)
			}).
		AddTextArea(labelAdditionalInformation, "", defaultFieldWidth, defaultFieldHeight, defaultMaxLength,
			func(text string) {
				m.Set(model.MetaKeyAdditionalInformation, text)
			}).
		AddButton(labelAdd, func() {
			d, err := os.ReadFile(path)
			if err != nil {
				c.NotifyPage(err.Error())
				return
			}

			s, err := os.Stat(path)
			if err != nil {
				c.NotifyPage(err.Error())
				return
			}

			if s.Size() > int64(storage.MaxFileSize) {
				c.NotifyPage(storage.ErrLargeFile.Error())
				return
			}
			f.Name = s.Name()
			f.Body = d

			_, err = c.storage.CreateItem(ctx, f, m)
			if err != nil {
				log.Error().Err(err).Msg("Item add error while add file")
				c.NotifyPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameAddFile)
			c.ItemsPage(ctx)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameAddFile)
		})

	form.
		SetTitle("Add file").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameAddFile, flex, true, true)
}

func (c *client) UpdateItemPage(ctx context.Context, item *storage.Item) {
	log.Printf("Invoked Update item page")

	switch t := item.Body.(type) {
	case *storage.Password:
		c.UpdatePasswordPage(ctx, item, t)
	case *storage.Card:
		c.UpdateCardPage(ctx, item, t)
	case *storage.Note:
		c.UpdateNotePage(ctx, item, t)
	case *storage.File:
		c.UpdateFilePage(ctx, item, t)
	default:
		log.Error().Msgf("Unknown item body type:%v", reflect.TypeOf(t))
	}
}

func (c *client) DeleteItemPage(ctx context.Context, i *storage.Item, name, itype string) {
	log.Printf("Invoked Delete item page, item(%v)", i.ID)

	text := fmt.Sprintf("Delete %s %s?", itype, name)

	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{buttonNameDelete}).
		AddButtons([]string{buttonNameCancel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == buttonNameCancel {
				c.pages.RemovePage(pageNameNotify)
				return
			}
			if buttonLabel == buttonNameDelete {
				err := c.storage.DeleteItem(ctx, i.ID)
				if err != nil {
					c.NotifyPage("error of delete item" + err.Error())
					return
				}

				c.pages.RemovePage(pageNameNotify)
				c.ItemsPage(ctx)
			}
		})

	c.pages.AddPage(pageNameNotify, modal, true, true)
}

func (c *client) StartSync(ctx context.Context) {
	log.Printf("Start sync")
	c.sync.Start(ctx)
}

func (c *client) StopSync(ctx context.Context) {
	c.sync.Stop(ctx)
}

func newTableCellTime(t time.Time) *tview.TableCell {
	nt := time.Time{}

	if t == nt {
		return tview.NewTableCell("")
	}

	return tview.NewTableCell(t.Local().Format(time.RFC3339)).SetTextColor(tcell.ColorWhite)
}
