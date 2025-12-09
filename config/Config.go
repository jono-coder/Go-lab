package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig
	DB   DBConfig
	Auth AuthConfig
}

type AppConfig struct {
	Env              string
	Protocol         string
	Host             string
	Port             uint16
	Root             string
	BaseUrl          string
	TimeoutInSeconds time.Duration
	Throttle         uint32
}

func (c *AppConfig) IsDev() bool {
	return c.Env == "dev"
}

type DBConfig struct {
	Driver string
	DSN    string
}

type AuthConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
}

func Load() (Config, error) {
	err := load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cfg := Config{
		App: AppConfig{
			Env:              getenv("APP_ENV", "dev"),
			Protocol:         getenv("APP_PROTOCOL", "http"),
			Host:             getenv("APP_HOST", "localhost"),
			Port:             uint16(getInt("APP_PORT", 8080)),
			Root:             getenv("APP_ROOT", "/"),
			TimeoutInSeconds: time.Second * time.Duration(getInt("APP_SERVICE_TIMEOUT", 5)),
			Throttle:         uint32(getInt("APP_THROTTLE", 10)),
		},
		DB: DBConfig{
			Driver: getenv("DB_DRIVER", "sqlite"),
			DSN:    getenvRequired("DB_DSN"),
		},
		Auth: AuthConfig{
			ClientID:     getenvRequired("AUTH_CLIENT_ID"),
			ClientSecret: getenvRequired("AUTH_CLIENT_SECRET"),
			TokenURL:     getenvRequired("AUTH_TOKEN_URL"),
		},
	}
	cfg.App.BaseUrl = fmt.Sprintf("%s://%s:%d%s", cfg.App.Protocol, cfg.App.Host, cfg.App.Port, cfg.App.Root)

	return cfg, nil
}

// Helpers //

func load() error {
	// Start at the working directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Walk up until we find .env or hit filesystem root
	for {
		try := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(try); statErr == nil {
			return godotenv.Load(try)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvRequired(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required environment variable: %s", key)
		return ""
	}
	return v
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid int for %s: %v", key, err)
	}
	return i
}

func splitComma(key string, def []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return strings.Split(v, ",")
}
