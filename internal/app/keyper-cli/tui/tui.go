package tui

import (
	"github.com/rivo/tview"
)

func NewTUI() *TUI {
	return &TUI{
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
}

type TUI struct {
	App   *tview.Application
	Pages *tview.Pages
}

func (t *TUI) LoginPage() *tview.Form {
	var form *tview.Form
	form = tview.NewForm().
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			// login
			_ = form.GetFormItem(0).(*tview.InputField).GetText()
			// password
			_ = form.GetFormItem(1).(*tview.InputField).GetText()

			var err error

			success, err := true, nil
			if err != nil || !success {
				modal := tview.NewModal().
					SetText("Login failed").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						t.Pages.RemovePage("error")
					})
				t.Pages.AddPage("error", modal, false, true)
			} else {
				t.Pages.SwitchToPage("secrets")
			}
		}).
		AddButton("Quit", func() {
			t.App.Stop()
		})
	form.SetBorder(true).SetTitle("Login").SetTitleAlign(tview.AlignCenter)
	return form
}

func (t *TUI) SecretsPage() *tview.List {
	var list *tview.List
	list = tview.NewList()
	list.AddItem("Error loading secrets", "", 0, nil)
	list.AddItem("Back", "Return to login", 0, func() {
		t.Pages.SwitchToPage("login")
	})
	list.SetBorder(true).SetTitle("Secrets").SetTitleAlign(tview.AlignCenter)
	return list
}

func ErrorPage() *tview.Form {
	return nil
}
