package config

type MailConfig struct {
	Username       string `envconfig:"MAIL_USERNAME"`
	Password       string `envconfig:"MAIL_PASSWORD"`
	Host           string `envconfig:"MAIL_SERVER"`
	Port           int    `envconfig:"MAIL_PORT"`
	From           string `envconfig:"MAIL_FROM"`
	FromName       string `envconfig:"MAIL_FROM_NAME"`
	StartTLS       bool   `envconfig:"MAIL_STARTTLS"`
	SSLTLS         bool   `envconfig:"MAIL_SSL_TLS"`
	UseCredentials bool   `envconfig:"USE_CREDENTIALS"`
	ValidateCerts  bool   `envconfig:"VALIDATE_CERTS"`
}
