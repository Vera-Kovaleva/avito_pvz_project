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

	// мне их нужно создать?
	"avito_pvz/internal/infra/container"
	"avito_pvz/internal/infra/log"
	//

	"github.com/gin-gonic/gin"
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
	os.Exit(Run(context.Background(), container.NewContainer()))
}

func Run(ctx context.Context, container *container.Container) int {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.ErrorContext(ctx, "Loading environment variables failed.", log.ErrorAttr(err))

		return exitDotEnvFailed
	}

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

	var eg errgroup.Group
	startHTTPServer(ctx, &eg, container.HTTPServer(ctx))
	if err := eg.Wait(); err != nil {
		slog.ErrorContext(ctx, "Runing servers failed.", log.ErrorAttr(err))

		return exitServersFailed
	}

	return exitOK
}

func startHTTPServer(ctx context.Context, eg *errgroup.Group, server *gin.Engine) {
	httpSrv := &http.Server{
		Addr:              os.Getenv("HTTP_ADDRESS"),
		Handler:           server.Handler(),
		ReadTimeout:       readimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	eg.Go(func() error {
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
