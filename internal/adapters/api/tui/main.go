package tui

import (
	"context"
	"fmt"

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
	pageNameError    = "error"

	buttonNameCancel = "Cancel"

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

func (c *client) RegisterPage() {
	log.Printf("Invoked Register Page")

	var (
		email           string
		password        string
		confirmPassword string
	)
	registerForm := tview.NewForm()
	registerForm.
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
				log.Printf("err:%v", err)
			}

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
	log.Printf("Invoked Login Page")

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
			//c.app.SetFocus(loginUnsuccessModal)
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
