package handlers

import (
	"errors"
	"net/http"

	serviceerrors "krizanauskas.github.com/mvp-proxy/internal/errors"

	"github.com/labstack/echo/v4"
	"krizanauskas.github.com/mvp-proxy/internal/services"
)

func ProxyHandler(c echo.Context) error {
	proxyService := services.NewProxyService(c)

	err := proxyService.ProxyRequest()
	var serviceErr *serviceerrors.ServiceError

	if err != nil {
		if errors.As(err, &serviceErr) {
			return proxyService.String(serviceErr.Code, serviceErr.Message)
		}

		return proxyService.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
