package app

import (
	"log/slog"

	"github.com/rivo/tview"
)

type App struct {
	log *slog.Logger

	tviewApp   *tview.Application
	tviewPages *tview.Pages

	shortcuts []rune
}

func NewApp(log *slog.Logger) *App {
	return &App{
		log:        log,
		tviewApp:   tview.NewApplication(),
		tviewPages: tview.NewPages(),
		shortcuts:  []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
	}
}

// Run runs the whole application.
func (a *App) Run() error {
	
	err := a.tviewApp.
		SetRoot(a.tviewPages, true).
		EnableMouse(true).
		EnablePaste(true).
		Run()

	if err != nil {
		return err
	}

	return nil
}
