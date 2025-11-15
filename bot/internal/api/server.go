package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"maxBot/internal/auth"
	"maxBot/internal/di"
)

// Server инкапсулирует HTTP-слой приложения.
type Server struct {
	cfg  RestConfig
	http *http.Server
}

// NewServer создаёт gin.Engine и HTTP-сервер.
func NewServer(cfg RestConfig, services *di.Services, validator *auth.Validator) (*Server, error) {
	if services == nil || services.EventService == nil {
		return nil, fmt.Errorf("event service is required")
	}
	if services.UserService == nil {
		return nil, fmt.Errorf("user service is required")
	}
	if validator == nil {
		return nil, fmt.Errorf("auth validator is required")
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	corsCfg := buildCORSConfig(cfg.AllowedOrigins)
	engine.Use(cors.New(corsCfg))

	apiV1 := engine.Group("/api/v1")
	authMW := newAuthMiddleware(validator)
	mapHandler := newMapHandler(services.EventService)
	mapHandler.register(apiV1, authMW)
	newUserHandler(services.UserService).register(apiV1, authMW)
	newAuthHandler(validator).register(apiV1)

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{
		cfg:  cfg,
		http: httpServer,
	}, nil
}

// Run запускает HTTPS сервер и блокирует выполнение.
func (s *Server) Run() error {
	return s.http.ListenAndServe()
	// return s.http.ListenAndServeTLS(s.cfg.TLSCertFile, s.cfg.TLSKeyFile)
}

// Shutdown останавливает сервер с graceful shutdown.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
