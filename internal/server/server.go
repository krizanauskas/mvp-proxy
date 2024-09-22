package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"krizanauskas.github.com/mvp-proxy/config/appconfig"
	"krizanauskas.github.com/mvp-proxy/internal/handlers"
	"krizanauskas.github.com/mvp-proxy/internal/middlewares"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg appconfig.ProxyServerConfig, proxyHandler handlers.ProxyHandler, m ...func(http.Handler) http.Handler) (*Server, error) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(cfg.MaxRequestDurationSec)*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		proxyHandler := http.HandlerFunc(proxyHandler.Handle)

		handler := middlewares.ChainMiddleware(proxyHandler, m...)

		handler.ServeHTTP(w, r)
	})

	httpServer := &http.Server{
		Addr:              cfg.Port,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 60,
		Handler:           handlerFunc,
	}

	server := &Server{
		httpServer,
	}

	return server, nil
}

func (s *Server) Start() error {
	fmt.Println("Starting server on", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
