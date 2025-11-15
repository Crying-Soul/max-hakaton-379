package auth

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	envJWTSecret  = "AUTH_JWT_SECRET"
	envAuthMaxAge = "AUTH_MAX_AGE"
	envSessionTTL = "AUTH_SESSION_TTL"
)

// LoadConfigFromEnv builds Config using env vars and provided bot token.
func LoadConfigFromEnv(botToken string) (Config, error) {
	cfg := Config{
		BotToken:   strings.TrimSpace(botToken),
		JWTSecret:  strings.TrimSpace(os.Getenv(envJWTSecret)),
		MaxAge:     parseDurationEnv(envAuthMaxAge),
		SessionTTL: parseDurationEnv(envSessionTTL),
	}

	if cfg.BotToken == "" {
		return Config{}, fmt.Errorf("bot token is required for auth config")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("%s is required", envJWTSecret)
	}

	return cfg, nil
}

func parseDurationEnv(key string) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return 0
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return d
}
