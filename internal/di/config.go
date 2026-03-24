package di

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Connection pooling constants (adjust as needed)
const (
	BunMaxOpenConns    = 25
	BunMaxIdleConns    = 5
	BunConnMaxLifetime = 5 * time.Minute
	BunConnMaxIdleTime = 5 * time.Minute
)

type Config struct {
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresSSLMode  string `env:"POSTGRES_SSLMODE" envDefault:"disable"` // disable, require, verify-full, etc.

	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"console"`

	Env  string `env:"ENV" envDefault:"development"`
	Port string `env:"PORT" envDefault:"8080"`
}

func NewConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := env.ParseAs[Config]()
	return &cfg, err
}
