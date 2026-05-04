package monitoring

import (
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func InitSentry(cfg config.SentryConfig, logger *zap.Logger) error {
	if !cfg.IsEnabled {
		logger.Info("Sentry is disabled")
		return nil
	}

	if cfg.DSN == "" {
		logger.Info("Sentry DSN not provided, skipping initialization")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:            cfg.DSN,
		Release:        cfg.Release,
		Environment:    cfg.Env,
		SendDefaultPII: cfg.SendDefaultPII,
		EnableTracing:  false,
	})

	if err != nil {
		return err
	}

	logger.Info("Sentry initialized for environment", zap.String("environment", cfg.Env))
	return nil
}

func CloseSentry(cfg config.SentryConfig, logger *zap.Logger) {
	if !cfg.IsEnabled {
		logger.Info("Sentry is disabled")
		return
	}
	sentry.Flush(2 * time.Second)
}
