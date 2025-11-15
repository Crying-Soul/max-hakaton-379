package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	godotenv "github.com/joho/godotenv"
	maxbot "github.com/rectid/max-bot-api-client-go"

	"maxBot/internal/api"
	"maxBot/internal/auth"
	"maxBot/internal/di"
	"maxBot/internal/repository"
)

// main запускает бота с graceful shutdown
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN for bot not set")
	}

	apiClient, err := maxbot.New(token)
	if err != nil {
		log.Fatalf("Failed to init API client: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	repo, err := repository.New(ctx, dsn, os.Getenv("MIGRATIONS_PATH"))
	if err != nil {
		log.Fatalf("Failed to init repository: %v", err)
	}
	defer repo.Close()

	services := di.NewServices(apiClient, repo)

	cfg, err := api.LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("HTTP config error: %v", err)
	}

	authCfg, err := auth.LoadConfigFromEnv(token)
	if err != nil {
		log.Fatalf("Auth config error: %v", err)
	}

	authValidator, err := auth.NewValidator(authCfg)
	if err != nil {
		log.Fatalf("Failed to init auth validator: %v", err)
	}

	server, err := api.NewServer(cfg, services, authValidator)
	if err != nil {
		log.Fatalf("Failed to init HTTP server: %v", err)
	}

	go func() {
		log.Printf("HTTPS server listening on %s", cfg.Addr)
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
			stop()
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown error: %v", err)
		}
	}()

	// bot, err := internal.NewBot(ctx, services)
	// if err != nil {
	// 	log.Fatalf("Failed to create bot: %v", err)
	// }
	// defer bot.Close()

	// bot.Start(ctx)

	// log.Println("Bot stopped")

	for {
	}
}
