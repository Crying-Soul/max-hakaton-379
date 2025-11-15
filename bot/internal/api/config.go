package api

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAddr            = ":8443"
	defaultReadTimeout     = 10 * time.Second
	defaultWriteTimeout    = 15 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

// RestConfig описывает настройки HTTP-сервера.
type RestConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	TLSCertFile     string
	TLSKeyFile      string
	AllowedOrigins  []string
}

// LoadConfigFromEnv читает настройки из переменных окружения.
func LoadConfigFromEnv() (RestConfig, error) {
	cfg := RestConfig{
		Addr:            valueOrDefault(os.Getenv("HTTP_ADDR"), defaultAddr),
		ReadTimeout:     durationOrDefault(os.Getenv("HTTP_READ_TIMEOUT"), defaultReadTimeout),
		WriteTimeout:    durationOrDefault(os.Getenv("HTTP_WRITE_TIMEOUT"), defaultWriteTimeout),
		ShutdownTimeout: durationOrDefault(os.Getenv("HTTP_SHUTDOWN_TIMEOUT"), defaultShutdownTimeout),
		TLSCertFile:     strings.TrimSpace(os.Getenv("HTTP_TLS_CERT_FILE")),
		TLSKeyFile:      strings.TrimSpace(os.Getenv("HTTP_TLS_KEY_FILE")),
		AllowedOrigins:  parseCSV(os.Getenv("HTTP_ALLOWED_ORIGINS")),
	}

	if cfg.TLSCertFile == "" || cfg.TLSKeyFile == "" {
		return RestConfig{}, fmt.Errorf("HTTP_TLS_CERT_FILE и HTTP_TLS_KEY_FILE обязательны для HTTPS")
	}

	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{
			"https://*.vercel.app",
			"https://maxapp.ru",
			"https://web.maxapp.ru",
			"https://messenger.max.ru",
		}
	}

	return cfg, nil
}

func valueOrDefault(val, def string) string {
	if strings.TrimSpace(val) == "" {
		return def
	}
	return val
}

func durationOrDefault(val string, def time.Duration) time.Duration {
	if strings.TrimSpace(val) == "" {
		return def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return def
	}
	return d
}

func parseCSV(val string) []string {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

// parseInt32Bound читает значение int32 из строки с ограничениями, возвращает дефолт при ошибке.
func parseInt32Bound(val string, min, max, def int32) int32 {
	if strings.TrimSpace(val) == "" {
		return def
	}
	n, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return def
	}
	res := int32(n)
	if res < min {
		return min
	}
	if res > max {
		return max
	}
	return res
}
