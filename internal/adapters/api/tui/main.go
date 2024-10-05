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

	pageNameError  = "error"
	pageNameNotify = "notify"

	buttonNameCancel = "Cancel"
	buttonNameOk     = "Ok"

	formPassword = "password"
	formCard     = "card"
	formNote     = "note"
	formFile     = "file"
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
	colDesc
	colCreated
	colModified
	colType
	colHash
	colVersion
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

			c.NotifyAndSwitch2Page("Success login", c.ItemsPage)
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

func (c *client) ItemsPage(ctx context.Context, offset int, limit int) {
	log.Printf("Invoked Items Page")
	curOst := offset
	curLt := limit

	list, err := c.grpc.ListItem(context.Background(), email, password)

	table := tview.NewTable().
		SetCell(0, colID, addTableHeaderCell("ID")).
		SetCell(0, colDesc, addTableHeaderCell(strings.ToUpper(fnDescription))).
		SetCell(0, colCreated, addTableHeaderCell("CREATED")).
		SetCell(0, colModified, addTableHeaderCell("MODIFIED")).
		SetCell(0, colType, addTableHeaderCell("TYPE")).
		SetCell(0, colHash, addTableHeaderCell("HASHSUM")).
		SetCell(0, colVersion, addTableHeaderCell("VERSION"))

	for r := curOst; r < len(list); r++ {
		record := rs[r]
		rn := r + 1

		table.SetCell(rn, colID, addTableCell(record.ID))
		table.SetCell(rn, colDesc, addTableHeaderCell(record.Description))
		table.SetCell(rn, colCreated, addTableHeaderCell(record.Created.Format(fnDateFormat)))
		table.SetCell(rn, colModified, addTableHeaderCell(record.Modified.Format(fnDateFormat)))
		table.SetCell(rn, colType, addTableHeaderCell(record.Type))
		table.SetCell(rn, colHash, addTableHeaderCell(record.Hashsum))
		table.SetCell(rn, colVersion, addTableHeaderCell(strconv.FormatInt(record.GetVersion(), 10)))

		if rn >= curOst+ui.recLimit {
			break
		}
	}
	table.SetSelectable(true, false)

	table.SetSelectedFunc(func(row int, column int) {
		recordID := table.GetCell(row, colID).Text
		if strings.TrimSpace(recordID) == "" {
			ui.displayErr("record id is empty")
			return
		}

		dataType := table.GetCell(row, colType).Text
		switch dataType {
		case string(models.AuthType):
			ui.displayUpdateAuth(ctx, recordID)
		case string(models.TextType):
			ui.displayUpdateText(ctx, recordID)
		case string(models.BinaryType):
			ui.displayUpdateBinary(ctx, recordID)
		case string(models.CardType):
			ui.displayUpdateCard(ctx, recordID)
		default:
			ui.displayErr("Unknow type")
		}
	})

	buttonsManageList := tview.NewForm().
		AddButton("<", func() {
			curOst := max(0, curOst-curLt)

			ui.pages.RemovePage(pageListRecords)
			ui.displayRecords(ctx, curOst, curLt)
		}).
		AddButton("Refresh", func() {
			ui.pages.RemovePage(pageListRecords)
			ui.displayRecords(ctx, curOst, curLt)
		}).
		AddButton(">", func() {
			curOst := curOst + curLt

			ui.pages.RemovePage(pageListRecords)
			ui.displayRecords(ctx, curOst, curLt)
		}).
		AddButton("Back to menu", func() {
			ui.pages.RemovePage(pageListRecords)
		})
	buttonsManageList.SetButtonsAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 0, 0)

	buttons := tview.NewForm().
		AddButton("Add auth", func() { ui.displayCreateAuth(ctx) }).
		AddButton("Add text", func() { ui.displayCreateText(ctx) }).
		AddButton("Add file", func() { ui.displayCreateBinary(ctx) }).
		AddButton("Add card", func() { ui.displayCreateCard(ctx) })

	buttons.SetButtonsAlign(tview.AlignLeft).SetBorderPadding(0, 0, 0, 0)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(buttons, 1, 1, false).
		AddItem(table, 0, 1, true).
		AddItem(buttonsManageList, 1, 1, false)

	flex.SetBorder(true)

	ui.pages.AddPage(pageListRecords, flex, true, true)

	c.pages.AddPage(pageNameItems, itemsFlexBox, true, true)
}

func (c *client) ItemsPage(ctx context.Context, offset, limit int) {
	curOst := offset
	curLt := limit

	rs, err := ui.authUser.GetRecords(ctx, ui.cache, curOst, curLt)
	if err != nil {
		ui.displayErr(fmt.Sprintf("an error occured while retrieving record list, err: %v", err))
		return
	}

	table := tview.NewTable()

	table.SetCell(0, colID, addTableHeaderCell("ID"))
	table.SetCell(0, colDesc, addTableHeaderCell(strings.ToUpper(fnDescription)))
	table.SetCell(0, colCreated, addTableHeaderCell("CREATED"))
	table.SetCell(0, colModified, addTableHeaderCell("MODIFIED"))
	table.SetCell(0, colType, addTableHeaderCell("TYPE"))
	table.SetCell(0, colHash, addTableHeaderCell("HASHSUM"))
	table.SetCell(0, colVersion, addTableHeaderCell("VERSION"))

	for r := curOst; r < len(rs); r++ {
		record := rs[r]
		rn := r + 1

		table.SetCell(rn, colID, addTableCell(record.ID))
		table.SetCell(rn, colDesc, addTableHeaderCell(record.Description))
		table.SetCell(rn, colCreated, addTableHeaderCell(record.Created.Format(fnDateFormat)))
		table.SetCell(rn, colModified, addTableHeaderCell(record.Modified.Format(fnDateFormat)))
		table.SetCell(rn, colType, addTableHeaderCell(record.Type))
		table.SetCell(rn, colHash, addTableHeaderCell(record.Hashsum))
		table.SetCell(rn, colVersion, addTableHeaderCell(strconv.FormatInt(record.GetVersion(), 10)))

		if rn >= curOst+ui.recLimit {
			break
		}
	}
	table.SetSelectable(true, false)

	table.SetSelectedFunc(func(row int, column int) {
		recordID := table.GetCell(row, colID).Text
		if strings.TrimSpace(recordID) == "" {
			ui.displayErr("record id is empty")
			return
		}

		dataType := table.GetCell(row, colType).Text
		switch dataType {
		case string(models.AuthType):
			ui.displayUpdateAuth(ctx, recordID)
		case string(models.TextType):
			ui.displayUpdateText(ctx, recordID)
		case string(models.BinaryType):
			ui.displayUpdateBinary(ctx, recordID)
		case string(models.CardType):
			ui.displayUpdateCard(ctx, recordID)
		default:
			ui.displayErr("Unknow type")
		}
	})

	buttonsManageList := tview.NewForm().
		AddButton("<", func() {
			curOst := max(0, curOst-curLt)

			ui.pages.RemovePage(pageListRecords)
			ui.displayRecords(ctx, curOst, curLt)
		}).
		AddButton("Refresh", func() {
			ui.pages.RemovePage(pageListRecords)
			ui.displayRecords(ctx, curOst, curLt)
		}).
		AddButton(">", func() {
			curOst := curOst + curLt

			ui.pages.RemovePage(pageListRecords)
			ui.displayRecords(ctx, curOst, curLt)
		}).
		AddButton("Back to menu", func() {
			ui.pages.RemovePage(pageListRecords)
		})
	buttonsManageList.SetButtonsAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 0, 0)

	buttons := tview.NewForm().
		AddButton("Add auth", func() { ui.displayCreateAuth(ctx) }).
		AddButton("Add text", func() { ui.displayCreateText(ctx) }).
		AddButton("Add file", func() { ui.displayCreateBinary(ctx) }).
		AddButton("Add card", func() { ui.displayCreateCard(ctx) })

	buttons.SetButtonsAlign(tview.AlignLeft).SetBorderPadding(0, 0, 0, 0)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(buttons, 1, 1, false).
		AddItem(table, 0, 1, true).
		AddItem(buttonsManageList, 1, 1, false)

	flex.SetBorder(true)

	ui.pages.AddPage(pageListRecords, flex, true, true)
}

func addTableHeaderCell(name string) *tview.TableCell {
	return tview.NewTableCell(name).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter)
}

func addTableCell(name string) *tview.TableCell {
	return tview.NewTableCell(name).
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignLeft)
}
