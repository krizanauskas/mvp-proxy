package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"krizanauskas.github.com/mvp-proxy/config/appconfig"
	"krizanauskas.github.com/mvp-proxy/internal/server"
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

	proxyServer, err := server.New(cfg.ProxyServer.Port)
	if err != nil {
		log.Fatalf("failed to init proxy server: %s", err.Error())
	}

	proxyServer.InitRoutes()

	if err := proxyServer.Start(); err != nil {
		log.Fatalf("failed to start proxy server: %s", err.Error())
	}
}
