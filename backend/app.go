// Package main provides the entry point for the nexus-tasks application.
package main

import (
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/cli"
	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/AshvinBambhaniya/nexus-tasks/logger"
	"github.com/AshvinBambhaniya/nexus-tasks/monitoring"
	"github.com/AshvinBambhaniya/nexus-tasks/routinewrapper"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func main() {
	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	logger, err := logger.NewRootLogger(cfg.Debug, cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	err = monitoring.InitSentry(cfg.Sentry, logger)
	if err != nil {
		logger.Error("Failed to initialize Sentry", zap.Error(err))
	}

	if cfg.Sentry.IsEnabled {
		// Sentry Go routine initialization
		sentryLoggedFunc := func() {
			err := recover()
			if err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 2)
			}
		}

		routinewrapper.Init(sentryLoggedFunc)
		defer sentryLoggedFunc()

		defer monitoring.CloseSentry(cfg.Sentry, logger)
	}

	err = cli.Init(cfg, logger)
	if err != nil {
		panic(err)
	}

}
