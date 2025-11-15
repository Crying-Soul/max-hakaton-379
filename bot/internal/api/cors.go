package api

import (
	"github.com/gin-contrib/cors"
)

func buildCORSConfig(allowed []string) cors.Config {
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"*"},
		AllowOrigins:     []string{"*"},
	}

	return cfg
}
