package services

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"krizanauskas.github.com/mvp-proxy/internal/errors"
)

type ProxyService struct {
	echo.Context
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func NewProxyService(c echo.Context) ProxyService {
	return ProxyService{
		c,
	}
}

func (ps ProxyService) ProxyRequest() error {
	return ps.handleHttpProxy()
}

func (ps ProxyService) handleHttpProxy() error {
	proxyReq, err := http.NewRequest(ps.Request().Method, ps.Request().URL.String(), ps.Request().Body)
	if err != nil {
		return &errors.ServiceError{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)}
	}
	proxyReq.Header = ps.Request().Header.Clone()

	proxyReq = proxyReq.WithContext(ps.Request().Context())

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return &errors.ServiceError{Code: http.StatusBadGateway, Message: http.StatusText(http.StatusBadGateway)}
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		if isHopByHopHeader(key) {
			continue
		}

		for _, value := range values {
			ps.Response().Header().Add(key, value)
		}
	}

	ps.Response().WriteHeader(resp.StatusCode)

	_, err = io.Copy(ps.Response(), resp.Body)
	if err != nil {
		// log error
	}

	return nil
}

func isHopByHopHeader(header string) bool {
	header = http.CanonicalHeaderKey(header)
	for _, h := range hopHeaders {
		if h == header {
			return true
		}
	}
	return false
}
