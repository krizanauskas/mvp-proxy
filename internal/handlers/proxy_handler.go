package handlers

import (
	"errors"
	"fmt"
	"net/http"

	serviceerrors "krizanauskas.github.com/mvp-proxy/internal/errors"
	"krizanauskas.github.com/mvp-proxy/internal/services"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("new request to %s \n", r.Host)

	proxyService := services.NewProxyService(w, r)

	err := proxyService.ProxyRequest()
	var serviceErr *serviceerrors.ServiceError

	if err != nil {
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), serviceErr.Code)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
