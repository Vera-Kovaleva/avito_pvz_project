package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/infra/database"
	"avito_pvz/internal/infra/log"
	"avito_pvz/internal/infra/noerr"
	"avito_pvz/internal/infra/repository"

	httpapi "avito_pvz/internal/adapters/http"
	oapi "avito_pvz/internal/generated/oapi"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const (
	exitOK = iota
	exitDotEnvFailed
	exitServersFailed
)

const (
	readimeout        = 100 * time.Millisecond
	readHeaderTimeout = 100 * time.Millisecond
)

func main() {
	os.Exit(Run(context.Background()))
}

func Run(ctx context.Context) int {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.ErrorContext(ctx, "Loading environment variables failed.", log.ErrorAttr(err))

		return exitDotEnvFailed
	}

	router := gin.Default()

	var stop context.CancelFunc
	ctx, stop = signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	switch gin.Mode() {
	case gin.DebugMode:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	provider := database.NewPostgresProvider(
		noerr.Must(pgxpool.New(ctx, os.Getenv("DB_CONNECTION"))),
	)
	defer provider.Close()

	pvzService := domain.NewPVZService(
		provider,
		repository.NewPVZ(),
		repository.NewProduct(),
		repository.NewReceptions(),
	)

	usersService := domain.NewUserService(
		provider,
		repository.NewUsers(),
		domain.HashPassword,
		domain.CompareHashAndPassword,
		domain.GenerateToken,
		domain.AuthenticateByToken,
	)

	receptionsService := domain.NewReceptionService(
		provider,
		repository.NewReceptions(),
		repository.NewProduct(),
	)

	oapi.RegisterHandlers(router, oapi.NewStrictHandler(httpapi.NewServer(pvzService, receptionsService, usersService), nil))

	var eg errgroup.Group
	startHTTPServer(ctx, &eg, router)
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "Runing servers failed.", log.ErrorAttr(err))

		return exitServersFailed
	}

	return exitOK
}

func startHTTPServer(ctx context.Context, eg *errgroup.Group, router *gin.Engine) {
	httpSrv := &http.Server{
		Addr:              os.Getenv("HTTP_ADDRESS"),
		Handler:           router,
		ReadTimeout:       readimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	eg.Go(func() error {
		slog.InfoContext(ctx, "Starting HTTP server", slog.String("addr", httpSrv.Addr))
		err := httpSrv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		return err
	})
	eg.Go(func() error {
		<-ctx.Done()

		return httpSrv.Shutdown(ctx)
	})
}
