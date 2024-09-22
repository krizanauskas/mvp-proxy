package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"krizanauskas.github.com/mvp-proxy/config/appconfig"
	"krizanauskas.github.com/mvp-proxy/internal/handlers"
	"krizanauskas.github.com/mvp-proxy/internal/middlewares"
	"krizanauskas.github.com/mvp-proxy/internal/server"
	"krizanauskas.github.com/mvp-proxy/internal/services"
	"krizanauskas.github.com/mvp-proxy/internal/storage"
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

	storage.Initialize(cfg.ProxyServer.AllowedDataMB)

	userHistoryStore := storage.NewUserHistoryStore()
	useBandwidthStore := storage.NewUserBandwidthStore()

	userService := services.NewUserService(userHistoryStore, useBandwidthStore)

	authService := services.NewAuthService()
	authMiddleware := middlewares.NewBasicAuthMiddleware(authService)
	bandwidthLimitMiddleware := middlewares.NewBandwidthLimitMiddleware(userService)

	proxyHandler := handlers.NewProxyHandler(userService)

	proxyServer, err := server.New(cfg.ProxyServer, proxyHandler, authMiddleware.Middleware, bandwidthLimitMiddleware.Middleware)
	if err != nil {
		log.Fatalf("failed to init proxy server: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	historyHandler := handlers.NewHistoryHandler(userHistoryStore)
	usageInfoHandler := handlers.NewUsageInfoHandler(useBandwidthStore)

	go startStatusServer(cfg.StatusServer, historyHandler, usageInfoHandler)

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

func startStatusServer(cfg appconfig.StatusServerConfig, historyHandler handlers.HistoryHandler, usageInfoHandler handlers.UsageInfoHandler) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/history", historyHandler.Handle)
	mux.HandleFunc("/usage-limits", usageInfoHandler.Handle)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}

	fmt.Printf("Starting status server on %s \n", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Health and history server error: %v", err)
	}
}
