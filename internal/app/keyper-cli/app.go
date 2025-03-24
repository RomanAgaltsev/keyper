package app

import (
	"log/slog"

	"github.com/rivo/tview"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-cli/tui"
	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
)

type App struct {
	log *slog.Logger

	tviewApp   *tview.Application
	tviewPages *tview.Pages
}

func NewApp(log *slog.Logger) *App {
	return &App{
		log:        log,
		tviewApp:   tview.NewApplication(),
		tviewPages: tview.NewPages(),
	}
}

// Run runs the whole application.
func (a *App) Run() error {
	a.tviewPages.AddPage("login", tui.LoginPage(), true, true)
	a.tviewPages.AddPage("secrets", tui.SecretsPage(), true, true)

	if err := a.tviewApp.SetRoot(a.tviewPages, true).Run(); err != nil {
		a.log.Error("running TUI application", sl.Err(err))
	}

	return nil
}
