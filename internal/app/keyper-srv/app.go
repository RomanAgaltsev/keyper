package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/server"
	"github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/service"
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

	pprofAddr := fmt.Sprintf("%s:%v", a.cfg.PPROF.Host, a.cfg.PPROF.Port)
	gatewayAddr := fmt.Sprintf("%s:%v", a.cfg.REST.Host, a.cfg.REST.Port)
	grpcAddr := fmt.Sprintf("%s:%v", a.cfg.GRPC.Host, a.cfg.GRPC.Port)

	/*
		Pprof server
	*/

	pprofServer := server.NewPprofServer(pprofAddr)

	g.Go(func() (err error) {
		const op = "app.RunPPROFServer"

		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
			}
		}()

		a.log.Info(op, "addr", pprofAddr)

		if err = pprofServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			a.log.Error("running pprof server", sl.Err(err))
			return fmt.Errorf("running pprof server has failed: %w", err)
		}

		a.log.Info(fmt.Sprintf("pprof started on %s", pprofAddr))

		return nil
	})

	/*
		Gateway server
	*/

	gatewayServer := server.NewGatewayServer(grpcAddr, gatewayAddr)

	g.Go(func() (err error) {
		const op = "app.RunGatewayServer"

		a.log.Info(op, "addr", gatewayAddr)

		if err = gatewayServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			a.log.Error(op, sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	/*
		gRPC server
	*/

	userService := service.NewUserService(a.cfg.App)
	secretService := service.NewSecretService(a.cfg.App)

	gRPCServer := server.NewGRPCServer(a.cfg.GRPC, userService, secretService)

	g.Go(func() (err error) {
		const op = "app.RunGRPCServer"

		listen, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			a.log.Error(op, sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
		defer func() { _ = listen.Close() }()

		a.log.Info(op, "addr", listen.Addr().String())

		if err = gRPCServer.Serve(listen); err != nil {
			a.log.Error(op, sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})

	/*
		Shutdown
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

	g.Go(func() error {
		defer a.log.Info("gateway server has been shut down")

		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()

		if err := gatewayServer.Shutdown(shutdownTimeoutCtx); err != nil {
			a.log.Error("gateway server shut down", sl.Err(err))
		}

		return nil
	})

	g.Go(func() error {
		defer a.log.Info("GRPC server has been shut down")

		<-ctx.Done()

		const op = "app.StopGRPCServer"

		a.log.With(slog.String("op", op)).
			Info("stopping gRPC server", slog.Int("port", a.cfg.GRPC.Port))

		// Используем встроенный в gRPCServer механизм graceful shutdown
		gRPCServer.GracefulStop()

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
