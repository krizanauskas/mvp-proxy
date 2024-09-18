package server

import (
	"fmt"
	"net/http"

	"krizanauskas.github.com/mvp-proxy/internal/handlers"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

func New(serverPort string) (*Server, error) {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr: serverPort,
	}

	server := &Server{
		httpServer,
		mux,
	}

	server.httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	s.mux.HandleFunc("/*", handlers.ProxyHandler)
}

func (s *Server) Start() error {
	fmt.Println("Starting server on", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
