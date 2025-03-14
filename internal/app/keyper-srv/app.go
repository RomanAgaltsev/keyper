package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/server"
	"github.com/RomanAgaltsev/keyper/internal/config"
	"github.com/RomanAgaltsev/keyper/internal/logger/sl"
)

const (
	timeoutServerShutdown = 5 * time.Second
	timeoutShutdown       = 10 * time.Second
)

// App struct of the application.
type App struct {
	cfg *config.Config
	log *slog.Logger
}

// NewApp creates new application.
func NewApp(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

// Run runs the whole application.
func (a *App) Run() error {
	appCtx, cancelAppCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelAppCtx()

	g, ctx := errgroup.WithContext(appCtx)

	context.AfterFunc(ctx, func() {
		ctx, cancelAfter := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelAfter()

		<-ctx.Done()
		a.log.Error("failed to gracefully shutdown the service")
	})

	/*
		Pprof server
	*/

	pprofServer := server.NewPprofServer(a.cfg.Pprof)

	g.Go(func() (err error) {
		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
			}
		}()

		a.log.Info(fmt.Sprintf("pprof started on %s", a.cfg.Pprof.Address))
		if err = pprofServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			a.log.Error("running pprof server", sl.Err(err))
			return fmt.Errorf("running pprof server has failed: %w", err)
		}
		return nil
	})

	/*
		Graceful shutdown
	*/

	g.Go(func() error {
		defer a.log.Info("pprof server has been shut down")

		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()

		if err := pprofServer.Shutdown(shutdownTimeoutCtx); err != nil {
			a.log.Error("pprof server shut down", sl.Err(err))
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
