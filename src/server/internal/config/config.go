package config

import (
	"os"
	"path/filepath"
	"time"
)

const developmentJWTSecret = "fraternity-development-secret-change-me"

type Config struct {
	Port              string
	DataFile          string
	JWTSecret         string
	AllowedOrigin     string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func Load() Config {
	return Config{
		Port:              envOrDefault("PORT", "4201"),
		DataFile:          envOrDefault("DATA_FILE", filepath.Join(os.TempDir(), "fraternity-users.json")),
		JWTSecret:         envOrDefault("JWT_SECRET", developmentJWTSecret),
		AllowedOrigin:     envOrDefault("ALLOWED_ORIGIN", "http://localhost:4200"),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
