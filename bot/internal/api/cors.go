package api

import (
	"strings"

	"github.com/gin-contrib/cors"
)

func buildCORSConfig(allowed []string) cors.Config {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Content-Length"},
	}

	var strict []string
	patterns := make([]string, 0)
	for _, origin := range allowed {
		if strings.Contains(origin, "*") {
			patterns = append(patterns, origin)
			continue
		}
		strict = append(strict, origin)
	}

	cfg.AllowOrigins = strict
	if len(patterns) > 0 {
		cfg.AllowOriginFunc = func(origin string) bool {
			for _, p := range patterns {
				if wildcardMatch(p, origin) {
					return true
				}
			}
			return false
		}
	}

	return cfg
}

// wildcardMatch проверяет, соответствует ли origin шаблону вида https://*.vercel.app
func wildcardMatch(pattern, origin string) bool {
	if pattern == origin {
		return true
	}
	if !strings.Contains(pattern, "*") {
		return false
	}
	parts := strings.Split(pattern, "*")
	if len(parts) != 2 {
		return false
	}
	prefix := parts[0]
	suffix := parts[1]
	return strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix)
}
