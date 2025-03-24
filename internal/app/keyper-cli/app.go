package app

import (
	"log/slog"

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
		tui: tui.NewTUI(),
	}
}

// Run runs the whole application.
func (a *App) Run() error {
	a.tui.Pages.AddPage("login", a.tui.LoginPage(), true, true)
	a.tui.Pages.AddPage("secrets", a.tui.SecretsPage(), true, true)

	if err := a.tui.App.SetRoot(a.tui.Pages, true).Run(); err != nil {
		a.log.Error("running TUI application", sl.Err(err))
	}

	return nil
}
