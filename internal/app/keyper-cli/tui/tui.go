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
	Register(ctx context.Context, user *model.User) (string, error)
	Login(ctx context.Context, user *model.User) (string, error)
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
		AddInputField("Server", "", 50, nil, nil).
		AddInputField("Login", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddTextArea("Token", "", 50, 5, 50, nil).
		AddButton("Register", func() {
			login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			user := model.User{
				Login:    login,
				Password: password,
			}

			// TODO: decide how to store token
			tokenString, err := t.user.Register(context.Background(), &user)
			if err != nil {
				modal := tview.NewModal().
					SetText("Register failed").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						t.Pages.RemovePage("error")
					})
				t.Pages.AddPage("error", modal, false, true)
			} else {
				form.GetFormItemByLabel("Token").(*tview.TextArea).SetText(tokenString, false)
			}
		}).
		AddButton("Login", func() {
			login := form.GetFormItemByLabel("Login").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

			user := model.User{
				Login:    login,
				Password: password,
			}

			// TODO: decide how to store token
			tokenString, err := t.user.Login(context.Background(), &user)
			if err != nil {
				modal := tview.NewModal().
					SetText("Register failed").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						t.Pages.RemovePage("error")
					})
				t.Pages.AddPage("error", modal, false, true)
			} else {
				form.GetFormItemByLabel("Token").(*tview.TextArea).SetText(tokenString, false)
			}
		}).
		AddButton("Quit", func() {
			t.App.Stop()
		})

	form.
		SetBorder(true).
		SetTitle("Keyper").
		SetTitleAlign(tview.AlignCenter)

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

func (t *TUI) SecretPage() *tview.Form {
	var form *tview.Form

	form = tview.NewForm().
		AddInputField("ID", "", 50, nil, nil).
		AddInputField("Name", "", 30, nil, nil).
		AddInputField("Type", "", 30, nil, nil).
		AddTextArea("Metadata", "", 50, 3, 50, nil).
		AddTextArea("Data", "", 50, 3, 50, nil).
		AddInputField("Comment", "", 50, nil, nil).
		AddInputField("Created", "", 30, nil, nil).
		AddInputField("Updated", "", 30, nil, nil).
		AddButton("Save", func() {}).
		AddButton("Copy", func() {}).
		AddButton("Back", func() {})

	form.
		SetBorder(true).
		SetTitle("Secret").
		SetTitleAlign(tview.AlignCenter)

	return form
}

func ErrorPage() *tview.Form {
	return nil
}
