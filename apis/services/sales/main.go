package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vishal2098govind/service/foundations/logger"
)

var build = "develop"
var buildDate = "YYYY-MM-DDTHH:MM:SSZ"

func main() {

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******** SEND ALERT *********")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	logger := logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events, build, buildDate)

	ctx := context.Background()

	if err := run(ctx, logger); err != nil {
		logger.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown
	log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
	defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}
