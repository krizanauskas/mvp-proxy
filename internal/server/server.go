package server

import (
	"github.com/labstack/echo/v4"
	"krizanauskas.github.com/mvp-proxy/internal/handlers"
)

type Server struct {
	*echo.Echo
}

func New() (*Server, error) {
	e := echo.New()

	s := &Server{
		e,
	}

	s.initRoutes()

	return s, nil
}

func (s *Server) initRoutes() {
	s.GET("/*", handlers.ProxyHandler)
}
