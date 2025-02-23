package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanAgaltsev/keyper/server/internal/app/keyper-srv/server"
	"github.com/RomanAgaltsev/keyper/server/internal/config"
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
func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	/*
		Pprof server
	*/

	pprofServer := server.NewPprofServer(a.cfg.Pprof)

	go func() {
		a.log.Info(fmt.Sprintf("pprof started on %s", a.cfg.Pprof.Address))
		if err := pprofServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Error("running pprof server: %s", err)
			cancel()
		}
	}()

	/*
		Graceful shutdown
	*/

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	a.log.Info("app is waiting for signal")

	select {
	case sig := <-quit:
		a.log.Info("received signal", "signal.Notify", sig.String())
		a.log.Info("server is shutting down")
	case <-ctx.Done():
		a.log.Info("received context done signal")
		a.log.Info("received context done signal", "ctx.Done", fmt.Sprintf("%v", ctx.Err()))
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pprofServer.Shutdown(ctx); err != nil {
		slog.Error("pprof shut down", err)
	} else {
		slog.Info("pprof shut down correctly")
	}

	return nil
}
