package tui

import (
	"context"
	"fmt"
	"strconv"

	grpc "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	pageNameWelcome  = "welcome"
	pageNameRegister = "register"
	pageNameLogin    = "login"
	pageNameItems    = "items"

	pageNameError  = "error"
	pageNameNotify = "notify"

	buttonNameCancel = "Cancel"
	buttonNameOk     = "Ok"

	formPassword = "password"
	formCard     = "card"
	formNote     = "note"
	formFile     = "file"

	pageSize = 40
)

const (
	fnDescription          = "Description"
	fnUsername             = "Username"
	fnPassword             = "Password"
	fnMetadata             = "Metadata"
	fnPath                 = "Path"
	fnText                 = "Text"
	fnNumber               = "Number"
	fnOwner                = "Owner"
	fnTerm                 = "Term"
	fnTemplateTermDesc     = "Template for term"
	fnTemplateHintTermDesc = "Please enter Term in format MM/YY, where MM - month, YY - year"
	fnDateFormat           = "02/01/2006 03:04.000"
)

const (
	colID = iota
	colType
	colName
	colDescription
	colCreateTime
	colUpdateTime
)

type client struct {
	grpc  grpc.GrpcClient
	app   *tview.Application
	pages *tview.Pages
}

func New(gc grpc.GrpcClient) *client {
	app := tview.NewApplication()
	pages := tview.NewPages()

	app.SetRoot(pages, true).EnableMouse(true)

	return &client{
		grpc:  gc,
		app:   app,
		pages: pages,
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

			c.NotifyAndSwitch2Page("Success login", c.ItemsPage
				func() {
					c.ItemsPage(context.Background(), 0, pageSize)
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
	log.Printf("Invoked Items Page, page_offset:%v, page_size:%v", page_offset, page_size)

	list, err := c.grpc.ListItems(context.Background(), page_offset, page_size)
	if err != nil {
		c.ErrorPage(fmt.Sprintf("error while get list items:%v", err))
		return
	}

	table := tview.NewTable().
		SetCell(0, colID, tview.NewTableCell("Id")).
		SetCell(0, colType, tview.NewTableCell("Type")).
		SetCell(0, colName, tview.NewTableCell("Name")).
		SetCell(0, colDescription, tview.NewTableCell("Description")).
		SetCell(0, colCreateTime, tview.NewTableCell("Create time")).
		SetCell(0, colUpdateTime, tview.NewTableCell("Update time"))

	for r, item := range list {
		table.SetCell(r, colID, tview.NewTableCell(strconv.FormatInt(item.ID, 10)))
		table.SetCell(r, colType, tview.NewTableCell(item.Type))
		table.SetCell(r, colName, tview.NewTableCell(item.Name))
	}

	table.SetSelectable(true, false)

	table.SetSelectedFunc(func(row, column int) {
		itemID := table.GetCell(row, colID).Text
		itemType := table.GetCell(row, colType).Text

		switch itemType {
		case string(grpc.ItemTypePassword):
			c.ViewPassword(ctx, itemID)
		case string(grpc.ItemTypeCard):
			c.ViewCard(ctx, itemID)
		case string(grpc.ItemTypeNote):
			c.ViewNote(ctx, itemID)
		case string(grpc.ItemTypeFile):
			c.ViewFile(ctx, itemID)
		}
	})

	navigateButtons := tview.NewForm().
		AddButton("Prev", func() {
			c.ItemsPage(ctx, max(page_offset-page_size), page_size)
		}).
		AddButton("Refresh", func() {
			c.ItemsPage(ctx, page_offset, page_size)
		}).
		AddButton("Next", func() {
			c.ItemsPage(ctx, page_offset+page_size, page_size)
		})

	navigateButtons.
		SetButtonsAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 0, 0)

	editButtons := tview.NewForm().
		AddButton("Add password", func() {
			//c.AddPassword(ctx)
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

	editButtons.
		SetButtonsAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 0, 0)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(editButtons, 1, 1, false).
		AddItem(table, 0, 1, true).
		AddItem(navigateButtons, 1, 1, false)

	flex.SetBorder(true)

	c.pages.AddPage(pageNameItems, flex, true, true)
}

func (c *client) ViewPassword(ctx context.Context, itemID string) {
	r, err := c.grpc.GetRecord(ctx, ui.cache, recordID)
	if err != nil {
		c.ErrorPage(err.Error())
	}

	a := &models.Auth{}
	if err := cbor.Unmarshal(r.Data, a); err != nil {
		ui.displayErr(fmt.Sprintf("unmarshal data auth binary, err: %v", err))
		return
	}

	desc := r.Description
	login := a.Login
	pass := a.Password
	metadata := r.Metadata
	mit := convertMetadataToString(r.Metadata)
	form := tview.NewForm().
		AddInputField(fnDescription, r.Description, defaultFieldWidth, nil, func(v string) {
			desc = v
		}).
		AddInputField(fnUsername, a.Login, defaultFieldWidth, nil, func(v string) {
			login = v
		}).
		AddInputField(fnPassword, a.Password, defaultFieldWidth, nil, func(v string) {
			pass = v
		}).
		AddTextArea(fnMetadata, mit, defaultFieldWidth, 0, 0, func(text string) {
			mit = text
		})

	form.SetTitle(pageUpdateAuthRecord).
		SetTitleAlign(tview.AlignLeft)

	buttons := tview.NewForm().
		AddButton(buttonUpdate, func() {
			auth := &models.Auth{
				Login:    login,
				Password: pass,
			}

			r, err := models.NewRecord(
				r.ID,
				desc,
				models.AuthType,
				r.Created,
				time.Now(),
				auth,
				metadata,
				false,
				r.Version,
			)
			if err != nil {
				ui.displayErr(err.Error())
				return
			}

			rows := splitMetadata(mit)
			m, err := models.NewMetadataFromStringArray(rows)
			if err != nil {
				ui.displayErr(err.Error())
				return
			}
			r.Metadata = m

			r.Version++
			_, err = ui.authUser.UpdateRecord(ctx, ui.cache, r)
			if err != nil {
				ui.displayErr(err.Error())
				return
			}

			ui.pages.RemovePage(pageUpdateAuthRecord)
		}).
		AddButton(buttonCancelDesc, func() { ui.pages.RemovePage(pageUpdateAuthRecord) })

	buttons.SetButtonsAlign(tview.AlignLeft).SetBorderPadding(0, 0, 0, 0)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).AddItem(buttons, 1, 1, false)

	ui.pages.AddPage(pageUpdateAuthRecord, flex, true, true)
}
