package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	LogLevel       string
	HTTPAddr       string
	DatabaseURL    string
	RedisAddr      string
	RequestTimeout time.Duration
}

func Load() (Config, error) {
	_ = godotenv.Load()
	timeout, err := getEnvDuration("REQUEST_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		AppEnv:         strings.ToLower(getEnv("APP_ENV", "development")),
		LogLevel:       strings.ToUpper(getEnv("LOG_LEVEL", "INFO")),
		HTTPAddr:       getEnv("HTTP_ADDR", "127.0.0.1:8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisAddr:      getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RequestTimeout: timeout,
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
func (c Config) validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.RedisAddr == "" {
		return errors.New("REDIS_ADDR is required")
	}
	if c.RequestTimeout <= 0 {
		return errors.New("REQUEST_TIMEOUT must be greater than 0")
	}
	if c.AppEnv != "development" && c.AppEnv != "production" {
		return fmt.Errorf("APP_ENV must be 'development' or 'production', got %q", c.AppEnv)
	}
	return nil
}
func getEnv(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s (%q): %w", key, raw, err)
	}
	return d, nil
}
