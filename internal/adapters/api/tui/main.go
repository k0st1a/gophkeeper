package tui

import (
	"fmt"

	grpc "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pageNameWelcome  = "welcome"
	pageNameRegister = "register"
	pageNameLogin    = "login"
	pageNameItems    = "items"

	formPassword = "password"
	formCard     = "card"
	formNote     = "note"
	formFile     = "file"
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
	welcomeList := tview.NewList().
		ShowSecondaryText(false).
		AddItem("Register", "", '1', func() {
			c.pages.RemovePage(pageNameWelcome)
			c.RegisterPage()
		}).
		AddItem("Login", "", '2', func() {
			c.pages.RemovePage(pageNameWelcome)
			c.LoginPage()
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

func (c *client) RegisterPage() {
	registerForm := tview.NewForm().
		AddInputField("Email", "", 30, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddPasswordField("Confirm password", "", 20, '*', nil).
		AddButton("Register", func() {
			c.pages.RemovePage(pageNameRegister)
			c.WelcomePage()
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

	// Login Page
	//loginUnsuccessModal := tview.NewModal().
	//	SetText("Unsuccess login: not implemeted. You will be returned on Welcome Page").
	//	AddButtons([]string{"Ok"})

	loginSuccessModal := tview.NewModal().
		SetText("Success login").
		AddButtons([]string{"Ok"})

	loginForm := tview.NewForm().
		AddInputField("Email", "", 30, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			//app.SetFocus(loginUnsuccessModal)
			c.pages.RemovePage(pageNameLogin)
			c.WelcomePage()
		}).
		AddButton("Back", func() {
			c.app.SetFocus(loginSuccessModal)
			c.app.SetFocus(c.pages)
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
