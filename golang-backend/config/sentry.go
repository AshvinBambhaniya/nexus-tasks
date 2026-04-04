package config

type SentryConfig struct {
	IsEnabled      bool   `envconfig:"SENTRY_ENABLED"`
	DSN            string `envconfig:"SENTRY_DSN"`
	Env            string `envconfig:"APP_ENV"`
	SendDefaultPII bool   `envconfig:"SENTRY_SEND_DEFAULT_PII"`
	Release        string `envconfig:"SENTRY_RELEASE"`
}
