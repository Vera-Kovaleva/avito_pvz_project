package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"avito_pvz/internal/domain"
	"avito_pvz/internal/infra/database"
	"avito_pvz/internal/infra/log"
	"avito_pvz/internal/infra/metrics"
	"avito_pvz/internal/infra/noerr"
	"avito_pvz/internal/infra/repository"

	httpapi "avito_pvz/internal/adapters/http"
	oapi "avito_pvz/internal/generated/oapi"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	metrics := metrics.NewMetrics()

	pvzService := domain.NewPVZService(
		provider,
		repository.NewPVZ(),
		repository.NewProduct(),
		repository.NewReceptions(),
		metrics,
	)

	usersService := domain.NewUserService(
		provider,
		repository.NewUsers(),
		metrics,
		domain.HashPassword,
		domain.CompareHashAndPassword,
		domain.GenerateToken,
		domain.AuthenticateByToken,
	)

	receptionsService := domain.NewReceptionService(
		provider,
		repository.NewReceptions(),
		repository.NewProduct(),
		metrics,
	)

	middlewares := []oapi.StrictMiddlewareFunc{
		func(f strictgin.StrictGinHandlerFunc, _ string) strictgin.StrictGinHandlerFunc {
			return func(ctx *gin.Context, request any) (any, error) {
				reqToken := ctx.Request.Header.Get("Authorization")
				if strings.HasPrefix(reqToken, "Bearer ") {
					token := strings.TrimPrefix(reqToken, "Bearer ")

					if user, err := domain.AuthenticateByToken(token); err == nil {
						ctx.Set(domain.CtxCurUserKey, user)
					}
				}

				return f(ctx, request)
			}
		},
		func(f strictgin.StrictGinHandlerFunc, _ string) strictgin.StrictGinHandlerFunc {
			return func(ctx *gin.Context, request any) (any, error) {
				start := time.Now()
				defer func() {
					metrics.IncRequests(time.Since(start))
				}()

				return f(ctx, request)
			}
		},
	}

	oapi.RegisterHandlers(
		router,
		oapi.NewStrictHandler(
			httpapi.NewServer(pvzService, receptionsService, usersService),
			middlewares,
		),
	)

	var eg errgroup.Group
	startHTTPServer(ctx, &eg, router)
	startPrometeusServer(ctx, &eg)

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

func startPrometeusServer(ctx context.Context, eg *errgroup.Group) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	metricsSrv := &http.Server{
		Addr:              ":9000",
		Handler:           mux,
		ReadTimeout:       readimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	eg.Go(func() error {
		slog.InfoContext(
			ctx,
			"Starting Prometheus metrics server",
			slog.String("addr", metricsSrv.Addr),
		)
		err := metricsSrv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		return err
	})
	eg.Go(func() error {
		<-ctx.Done()

		return metricsSrv.Shutdown(ctx)
	})
}
