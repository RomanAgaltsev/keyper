package tui

import (
	"context"

	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/service"
	"github.com/RomanAgaltsev/keyper/internal/model"
)

var (
	_ UserService   = (*service.UserService)(nil)
	_ SecretService = (*service.SecretService)(nil)
)

type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) error
}

type SecretService interface {
	Create(ctx context.Context, secret *model.Secret) (uuid.UUID, error)
	Update(ctx context.Context, userID uuid.UUID, secret *model.Secret) error
	UpdateData(ctx context.Context, userID uuid.UUID, secret *model.Secret) error
	Get(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error)
	GetData(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) (*model.Secret, error)
	List(ctx context.Context, userID uuid.UUID) (model.Secrets, error)
	Delete(ctx context.Context, userID uuid.UUID, secretID uuid.UUID) error
}

func NewTUI(user UserService, secret SecretService) *TUI {
	return &TUI{
		App:    tview.NewApplication(),
		Pages:  tview.NewPages(),
		user:   user,
		secret: secret,
	}
}

type TUI struct {
	App   *tview.Application
	Pages *tview.Pages

	user   UserService
	secret SecretService
}

func (t *TUI) LoginPage() *tview.Form {
	var form *tview.Form
	form = tview.NewForm().
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			//
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
