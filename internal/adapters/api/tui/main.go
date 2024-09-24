package tui

import (
	"fmt"

	grpc "github.com/k0st1a/gophkeeper/internal/adapters/api/grpc/client"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type client struct {
	grpc  grpc.GrpcClient
	app   *tview.Application
	pages *tview.Pages
}

func New(gc grpc.GrpcClient) *client {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Welcome Page
	welcomeList := tview.NewList().
		ShowSecondaryText(false).
		AddItem("Register", "", '1', func() {
			pages.SwitchToPage("register")
		}).
		AddItem("Login", "", '2', func() {
			pages.SwitchToPage("login")
		}).
		AddItem("Quit", "", 'q', func() {
			app.Stop()
		})

	welcomeList.
		SetTitle("Welcome").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	welcomeFlexBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(welcomeList, 0, 1, true)

	pages.AddPage("welcome", welcomeFlexBox, true, true)

	// Register Page
	registerForm := tview.NewForm().
		AddInputField("Email", "", 30, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddPasswordField("Confirm password", "", 20, '*', nil).
		AddButton("Register", func() {
			pages.SwitchToPage("welcome")
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("welcome")
		})

	registerForm.
		SetTitle("Enter register data").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	registerFlexBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(registerForm, 0, 1, true)

	pages.AddPage("register", registerFlexBox, true, false)

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
			pages.SwitchToPage("welcome")
		}).
		AddButton("Back", func() {
			app.SetFocus(loginSuccessModal)
			app.SetFocus(pages)
			pages.SwitchToPage("welcome")
		})

	loginForm.
		SetTitle("Enter login data").
		SetBorder(true).
		SetBorderColor(tcell.ColorSteelBlue)

	loginFlexBox := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(loginForm, 0, 1, true)

	pages.AddPage("login", loginFlexBox, true, false)

	// Items Page

	// Password Form

	// Card Form

	// Note Form

	// File Form

	app.SetRoot(pages, true).EnableMouse(true)

	return &client{
		app:   app,
		pages: pages,
	}
}

func (c *client) Run() error {
	err := c.app.Run()
	if err != nil {
		return fmt.Errorf("error of run tui client:%w", err)
	}

	return nil
}

func (c *client) Shutdown() error {
	return nil
}
