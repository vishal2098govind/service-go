package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/vishal2098govind/service/apis/services/auth/mux"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/keystore"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

var build = "develop"
var buildDate = "YYYY-MM-DDTHH:MM:SSZ"

func main() {

	traceIDFunc := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "************* SEND ALERT ****************")
		},
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelDebug, "AUTH", traceIDFunc, events, build, buildDate)

	ctx := context.Background()

	err := run(ctx, log)
	if err != nil {
		log.Error(ctx, "startup", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3001"`
			DebugHost          string        `conf:"default:0.0.0.0:3011"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
		Auth struct {
			Issuer    string `conf:"default:services-go service"`
			KeyPath   string `conf:"default:../../../zarf/keys/"`
			ActiveKid string `conf:"default:8234c8c5-0508-4301-bfdf-da1d515399a1"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Auth",
		},
	}

	const prefix = "AUTH"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if err == conf.ErrHelpWanted {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("conf parsing error: %w", err)
	}

	log.Info(ctx, "starting auth service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown completed")

	s, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", s)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ks := keystore.New()
	ks.LoadKeys(os.DirFS(cfg.Auth.KeyPath))

	auth, err := auth.New(auth.Config{
		Issuer:    cfg.Auth.Issuer,
		KeyLookup: &ks,
		Log:       log,
		ActiveKid: cfg.Auth.Issuer,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize auth: %w", err)
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(shutdown, log, auth),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverError := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverError <- api.ListenAndServe()
	}()

	select {
	case <-serverError:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("server couldn't stop down gracefully: %w", err)
		}

	}

	return nil
}
