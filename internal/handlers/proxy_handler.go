package handlers

import (
	"errors"
	"net/http"

	internalerrors "krizanauskas.github.com/mvp-proxy/internal/errors"
	"krizanauskas.github.com/mvp-proxy/internal/services"
)

type ProxyHandler struct {
	userService services.UserServiceI
}

func NewProxyHandler(userService services.UserServiceI) ProxyHandler {
	return ProxyHandler{
		userService,
	}
}

func (s ProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	proxyService := services.NewProxyService(w, r, s.userService)

	err := proxyService.ProxyRequest()
	var serviceErr *internalerrors.ServiceError

	if err != nil {
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), serviceErr.Code)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
