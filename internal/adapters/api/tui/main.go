package tui

import (
	"context"
	"fmt"
	"strconv"
	"time"

	grpc "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"
	"github.com/k0st1a/gophkeeper/internal/adapters/storage/inmemory"
	"github.com/k0st1a/gophkeeper/internal/pkg/client/model/item"
	"github.com/k0st1a/gophkeeper/internal/pkg/client/model/password"
	"github.com/k0st1a/gophkeeper/internal/ports"

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

	pageNameError  = "error"
	pageNameNotify = "notify"

	pageNameUpdatePassword = "update password"
	pageNameAddPassword    = "add password"

	// Имена кнопок.
	buttonNameCancel = "Cancel"
	buttonNameOk     = "Ok"
	buttonNameUpdate = "Update"

	formPassword = "password"
	formCard     = "card"
	formNote     = "note"
	formFile     = "file"

	labelItemName        = "Item name"
	labelItemDescription = "Item description"

	labelUserName = "User name"
	labelPassword = "Password"

	defaultFieldWidth = 30
)

const (
	columnID = iota
	columnName
	columnType
	columnDescription
	columnMarkDelete
	columnCreateTime
	columnUpdateTime
	columnUploadTime
	columnDownloadTime
)

type client struct {
	grpc    grpc.GrpcClient
	storage *inmemory.Storage
	app     *tview.Application
	pages   *tview.Pages
}

func New(gc grpc.GrpcClient, s *inmemory.Storage) *client {
	app := tview.NewApplication()
	pages := tview.NewPages()

	app.SetRoot(pages, true).EnableMouse(true)

	return &client{
		grpc:    gc,
		storage: s,
		app:     app,
		pages:   pages,
	}
}

func (c *client) Run() error {
	c.WelcomePage()
	err := c.app.Run()
	if err != nil {
		return fmt.Errorf("error of run tui client:%w", err)
	}

	return nil
}

func (c *client) Shutdown() error {
	return nil
}

func (c *client) WelcomePage() {
	log.Printf("Invoked Welcome Page")

	welcomeList := tview.NewList().
		ShowSecondaryText(false).
		AddItem("Login", "", '1', func() {
			c.pages.RemovePage(pageNameWelcome)
			c.LoginPage()
		}).
		AddItem("Register", "", '2', func() {
			c.pages.RemovePage(pageNameWelcome)
			c.RegisterPage()
		}).
		AddItem("Quit", "", 'q', func() {
			c.app.Stop()
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

func (c *client) ErrorPage(text string) {
	log.Printf("Invoked Error Page, text:%v", text)
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{buttonNameCancel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == buttonNameCancel {
				c.pages.RemovePage(pageNameError)
			}
		})

	c.pages.AddPage(pageNameError, modal, true, true)
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

func (c *client) RegisterPage() {
	log.Printf("Invoked Register Page")

	var (
		email           string
		password        string
		confirmPassword string
	)
	registerForm := tview.NewForm().
		AddInputField("Email", "", 30, nil, func(text string) {
			email = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Confirm password", "", 20, '*', func(text string) {
			confirmPassword = text
		}).
		AddButton("Register", func() {
			if email == "" {
				c.ErrorPage("Email is empty")
				return
			}

			if password == "" {
				c.ErrorPage("Password is empty")
				return
			}

			if confirmPassword == "" {
				c.ErrorPage("Confirm password is empty")
				return
			}

			if password != confirmPassword {
				c.ErrorPage("Passwords not equals")
				return
			}

			err := c.grpc.RegisterUser(context.Background(), email, password)
			if err != nil {
				c.ErrorPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameRegister)

			c.NotifyAndSwitch2Page("Success register", c.WelcomePage)
		}).
		AddButton("Cancel", func() {
			c.pages.RemovePage(pageNameRegister)
			c.WelcomePage()
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

func (c *client) LoginPage() {
	log.Printf("Invoked Login Page")

	var (
		email    string
		password string
	)
	loginForm := tview.NewForm().
		AddInputField("Email", "", 30, nil, func(text string) {
			email = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			if email == "" {
				c.ErrorPage("Email is empty")
				return
			}

			if password == "" {
				c.ErrorPage("Password is empty")
				return
			}

			err := c.grpc.LoginUser(context.Background(), email, password)
			if err != nil {
				c.ErrorPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameLogin)
			//c.ItemsPage()
			log.Printf("Success login fast")

			c.NotifyAndSwitch2Page("Success login",
				func() {
					c.pages.RemovePage(pageNameLogin)
					c.addTempPasswords()
					c.ItemsPage(context.Background(), 0, 20)
				})
			log.Printf("Success login fast")
		}).
		AddButton("Cancel", func() {
			c.pages.RemovePage(pageNameLogin)
			c.WelcomePage()
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

func (c *client) ItemsPage(ctx context.Context, page_offset, page_size int32) {
	log.Printf("Invoked Items Page")

	list := c.storage.ListItems(context.Background())

	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetSeparator(' ').
		SetCell(0, columnID, tview.NewTableCell("Id").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnName, tview.NewTableCell("Name").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnType, tview.NewTableCell("Type").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnDescription, tview.NewTableCell("Description").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnMarkDelete, tview.NewTableCell("Marked for delete?").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnCreateTime, tview.NewTableCell("Create time").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnUpdateTime, tview.NewTableCell("Update time").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnUploadTime, tview.NewTableCell("Upload time").SetSelectable(false).SetTextColor(tcell.ColorYellow)).
		SetCell(0, columnDownloadTime, tview.NewTableCell("Download time").SetSelectable(false).SetTextColor(tcell.ColorYellow))

	table.SetBorder(true).SetTitle("Table")

	for i, item := range list {
		row := i + 1
		table.
			SetCell(row, columnID, tview.NewTableCell(item.ID).SetSelectable(false).SetTextColor(tcell.ColorDarkCyan)).
			SetCell(row, columnName, tview.NewTableCell(item.Name).SetTextColor(tcell.ColorWhite)).
			SetCell(row, columnType, tview.NewTableCell(item.Type).SetTextColor(tcell.ColorWhite)).
			SetCell(row, columnDescription, tview.NewTableCell(item.Description).SetTextColor(tcell.ColorWhite)).
			SetCell(row, columnMarkDelete, tview.NewTableCell(strconv.FormatBool(item.MarkDelete)).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignRight)).
			SetCell(row, columnUploadTime, newTableCellTime(item.UploadTime).SetSelectable(false)).
			SetCell(row, columnCreateTime, newTableCellTime(item.CreateTime).SetSelectable(false)).
			SetCell(row, columnUpdateTime, newTableCellTime(item.UpdateTime).SetSelectable(false)).
			SetCell(row, columnDownloadTime, newTableCellTime(item.DownloadTime).SetSelectable(false))
	}

	table.SetSelectedFunc(func(row, column int) {
		itemName := table.GetCell(row, columnName).Text
		itemType := table.GetCell(row, columnType).Text

		switch itemType {
		case string(ports.ItemTypePassword):
			c.UpdatePasswordPage(ctx, itemName)
		case string(ports.ItemTypeCard):
			//c.ViewCard(ctx, itemName)
		case string(ports.ItemTypeNote):
			//c.ViewNote(ctx, itemName)
		case string(ports.ItemTypeFile):
			//c.ViewFile(ctx, itemName)
		}
	})

	table.
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	buttons := tview.NewForm().
		AddButton("Refresh", func() {
			c.ItemsPage(ctx, 0, 0)
		}).
		AddButton("Add password", func() {
			c.AddPasswordPage(ctx)
		}).
		AddButton("Add card", func() {
			//c.AddCard(ctx)
		}).
		AddButton("Add notes", func() {
			//c.AddNotes(ctx)
		}).
		AddButton("Add file", func() {
			//c.AddFile(ctx)
		})

	buttons.
		SetButtonsAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 0, 0)

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				AddItem(buttons, 1, 1, false).
				AddItem(table, 0, 1, true).
				SetDirection(tview.FlexRow), 0, 1, true)

	flex.
		SetTitle("Items page").
		SetBorder(true)

	c.pages.AddPage(pageNameItems, flex, true, true)
}

func (c *client) UpdatePasswordPage(ctx context.Context, name string) {
	i, err := c.storage.GetItem(ctx, name)
	if err != nil {
		c.ErrorPage(err.Error())
	}

	p, err := password.Deserialize(i.Data)
	if err != nil {
		c.ErrorPage(err.Error())
		return
	}

	form := tview.NewForm().
		AddInputField(labelItemName, i.Name, defaultFieldWidth, nil, func(text string) {
			i.Name = text
		}).
		AddInputField(labelItemDescription, i.Description, defaultFieldWidth, nil, func(text string) {
			i.Description = text
		}).
		AddInputField(labelUserName, p.UserName, defaultFieldWidth, nil, func(text string) {
			p.UserName = text
		}).
		AddInputField(labelPassword, p.Password, defaultFieldWidth, nil, func(text string) {
			p.Password = text
		}).
		AddButton(buttonNameUpdate, func() {
			d, err := password.Serialize(p)
			if err != nil {
				log.Printf("Password serialize error while update password:%w", err)
				c.ErrorPage(err.Error())
				return
			}

			log.Printf("Password serialized")

			i.Data = d

			err = c.storage.UpdateItem(ctx, i)
			if err != nil {
				log.Printf("Item update error while update password:%w", err)
				c.ErrorPage(err.Error())
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
	i := item.New()
	i.Type = ports.ItemTypePassword

	p := &password.Password{}

	form := tview.NewForm().
		AddInputField(labelItemName, i.Name, defaultFieldWidth, nil, func(text string) {
			i.Name = text
		}).
		AddInputField(labelItemDescription, i.Description, defaultFieldWidth, nil, func(text string) {
			i.Description = text
		}).
		AddInputField(labelUserName, p.UserName, defaultFieldWidth, nil, func(text string) {
			p.UserName = text
		}).
		AddInputField(labelPassword, p.Password, defaultFieldWidth, nil, func(text string) {
			p.Password = text
		}).
		AddButton(buttonNameOk, func() {
			d, err := password.Serialize(p)
			if err != nil {
				log.Printf("Password serialize error while add password:%w", err)
				c.ErrorPage(err.Error())
				return
			}

			log.Printf("Password serialized")

			i.Data = d

			err = c.storage.AddItem(ctx, i)
			if err != nil {
				log.Printf("Item add error while add password:%w", err)
				c.ErrorPage(err.Error())
				return
			}

			c.pages.RemovePage(pageNameAddPassword)
			c.ItemsPage(ctx, 0, 0)
		}).
		AddButton(buttonNameCancel, func() {
			c.pages.RemovePage(pageNameAddPassword)
		})
		//SetButtonsAlign(tview.AlignLeft).
		//SetBorderPadding(0, 0, 0, 0)

	form.
		SetTitle("Add password").
		SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true)

	c.pages.AddPage(pageNameAddPassword, flex, true, true)
}

func (c *client) addTempPasswords() {
	i := item.New()
	i.Type = ports.ItemTypePassword
	i.Name = "Name 1"
	i.Description = "Description 1"

	p := &password.Password{UserName: "UserName 1", Password: "Password 2"}
	d, _ := password.Serialize(p)
	i.Data = d

	_ = c.storage.AddItem(context.Background(), i)

	i = item.New()
	i.Type = ports.ItemTypePassword
	i.Name = "Name 2"
	i.Description = "Description 2"

	p = &password.Password{UserName: "UserName 2", Password: "Password 3"}
	d, _ = password.Serialize(p)
	i.Data = d

	_ = c.storage.AddItem(context.Background(), i)

	i = item.New()
	i.Type = ports.ItemTypePassword
	i.Name = "Name 3"
	i.Description = "Description 3"

	p = &password.Password{UserName: "UserName 3", Password: "Password 4"}
	d, _ = password.Serialize(p)
	i.Data = d

	_ = c.storage.AddItem(context.Background(), i)
}

func newTableCellTime(t time.Time) *tview.TableCell {
	nt := time.Time{}

	if t == nt {
		return tview.NewTableCell("")
	}

	return tview.NewTableCell(t.Format(time.RFC3339)).SetTextColor(tcell.ColorWhite)
}
