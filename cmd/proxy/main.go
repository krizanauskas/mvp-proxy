package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"krizanauskas.github.com/mvp-proxy/config/appconfig"
	"krizanauskas.github.com/mvp-proxy/internal/middlewares"
	"krizanauskas.github.com/mvp-proxy/internal/server"
	"krizanauskas.github.com/mvp-proxy/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to init .env: %s ", err.Error())
	}

	env := os.Getenv("APP_ENV")

	cfg, err := appconfig.Init(env)
	if err != nil {
		log.Fatalf("failed to init config: %s", err.Error())
	}

	authService := services.NewAuthService()
	authMiddleware := middlewares.NewBasicAuthMiddleware(authService)

	proxyServer, err := server.New(cfg.ProxyServer, authMiddleware.Middleware)
	if err != nil {
		log.Fatalf("failed to init proxy server: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := proxyServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start proxy server: %s", err.Error())
		}
	}()

	<-quit
	log.Println("Shutting down proxy server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := proxyServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
	os.Exit(0)
}
