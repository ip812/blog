package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Local Environment = "local"
	Prod  Environment = "prod"
)

func (e Environment) IsValid() bool {
	switch e {
	case Local, Prod:
		return true
	default:
		return false
	}
}

type Config struct {
	App struct {
		Env    Environment
		Domain string
		Port   string
	}

	Database struct {
		Name     string
		Endpoint string
		SSLMode  string
		Username string
		Password string
	}
}

func New() *Config {
	cfg := &Config{}

	cfg.App.Env = env
	if !cfg.App.Env.IsValid() {
		cfg.App.Env = Prod
	}

	err := godotenv.Load()
	if cfg.App.Env == Local && err != nil {
		panic(err)
	}

	cfg.App.Domain = os.Getenv("APP_DOMAIN")
	cfg.App.Port = os.Getenv("APP_PORT")
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.Endpoint = os.Getenv("DB_ENDPOINT")
	cfg.Database.SSLMode = os.Getenv("DB_SSL_MODE")
	cfg.Database.Username = os.Getenv("DB_USERNAME")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")

	return cfg
}
