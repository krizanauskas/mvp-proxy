package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"krizanauskas.github.com/mvp-proxy/config/appconfig"
	"krizanauskas.github.com/mvp-proxy/internal/handlers"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

func New(cfg appconfig.ProxyServerConfig) (*Server, error) {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:              cfg.Port,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 60,
	}

	server := &Server{
		httpServer,
		mux,
	}

	server.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(cfg.MaxRequestDurationSec)*time.Second)
		defer cancel()

		r = r.WithContext(ctx)

		if r.Method == http.MethodConnect {
			handlers.ProxyHandler(w, r)
		} else {
			// Forward other requests to the ServeMux for routing
			server.mux.ServeHTTP(w, r)
		}
	})

	return server, nil
}

func (s *Server) InitRoutes() {
	s.mux.HandleFunc("/", handlers.ProxyHandler)
}

func (s *Server) Start() error {
	fmt.Println("Starting server on", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
