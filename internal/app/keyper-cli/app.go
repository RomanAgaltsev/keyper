package app

import (
	"log/slog"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/client"
	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/service"
	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/tui"
	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
)

type App struct {
	log *slog.Logger
	tui *tui.TUI
}

func NewApp(log *slog.Logger) *App {
	return &App{
		log: log,
	}
}

// Run runs the whole application.
func (a *App) Run() error {
	userClient := client.NewUserClient()
	userService := service.NewUserService(a.log, userClient)

	secretClient := client.NewSecretClient()
	secretService := service.NewSecretService(a.log, secretClient)

	a.tui = tui.NewTUI(userService, secretService)

	a.tui.Pages.AddPage("login", a.tui.LoginPage(), true, false)
	a.tui.Pages.AddPage("secrets", a.tui.SecretsPage(), true, false)
	a.tui.Pages.AddPage("secret", a.tui.SecretPage(), true, true)

	if err := a.tui.App.SetRoot(a.tui.Pages, true).Run(); err != nil {
		a.log.Error("running TUI application", sl.Err(err))
	}

	return nil
}
