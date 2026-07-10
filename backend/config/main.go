package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// AppConfig type AppConfig
type AppConfig struct {
	IsDevelopment      bool   `envconfig:"IS_DEVELOPMENT"`
	Debug              bool   `envconfig:"DEBUG"`
	Env                string `envconfig:"APP_ENV"`
	Port               string `envconfig:"APP_PORT"`
	Secret             string `envconfig:"JWT_SECRET"`
	JwtExpirationHours int    `envconfig:"JWT_EXPIRATION_HOURS" default:"24"`
	AllowedOrigins     string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:3000"`
	MQ                 MQConfig
	DB                 DBConfig
	Sentry             SentryConfig
	Mail               MailConfig
	OpenAIKey          string `envconfig:"OPENAI_API_KEY"`
	OpenAIBaseURL      string `envconfig:"OPENAI_BASE_URL"`
	OpenAIModel        string `envconfig:"OPENAI_MODEL" default:"gpt-4o-mini"`
}

// GetConfig Collects all configs
func GetConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning .env file not found, scanning from OS ENV")
	}

	AllConfig = AppConfig{}

	err = envconfig.Process("", &AllConfig)
	if err != nil {
		log.Fatal(err)
	}

	return AllConfig
}

// GetConfigByName Collects all configs
func GetConfigByName(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

// LoadTestEnv loads environment variables from .env.testing file
func LoadTestEnv() AppConfig {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(fmt.Sprintf("%s/.env.testing", cwd))
	if err != nil {
		log.Fatal(err)
	}
	return GetConfig()
}
