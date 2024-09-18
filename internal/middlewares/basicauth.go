package middlewares

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"krizanauskas.github.com/mvp-proxy/internal/services"
)

const ProxyUsername = "user"
const ProxyPassword = "pass"

type credentials struct {
	username string
	password string
}

type AuthMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

type BasicAuthMiddleware struct {
	authService services.AuthServiceI
}

func NewBasicAuthMiddleware(authService services.AuthServiceI) BasicAuthMiddleware {
	return BasicAuthMiddleware{
		authService,
	}
}

func (m BasicAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials, err := m.parseCredentials(r)
		if err != nil {
			w.Header().Set("Proxy-Authenticate", `Basic realm="Proxy"`)
			http.Error(w, "failed to parse credentials", http.StatusProxyAuthRequired)

			fmt.Printf("error: %s", err.Error())
			return
		}

		if !m.authService.Authenticate(credentials.username, credentials.password) {
			w.Header().Set("Proxy-Authenticate", `Basic realm="Proxy"`)
			http.Error(w, "failed to authenticate", http.StatusProxyAuthRequired)

			fmt.Printf("failed to authenticate")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m BasicAuthMiddleware) parseCredentials(r *http.Request) (credentials, error) {
	authHeader := r.Header.Get("Proxy-Authorization")
	if authHeader == "" {
		return credentials{}, fmt.Errorf("empty credentials")
	}

	if !strings.HasPrefix(authHeader, "Basic ") {
		return credentials{}, fmt.Errorf("no 'Basic' prefix")
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
	if err != nil {
		return credentials{}, fmt.Errorf("failed to decode string: %s", err.Error())
	}

	// Split into username and password
	creds := strings.SplitN(string(payload), ":", 2)
	if len(creds) != 2 {
		return credentials{}, fmt.Errorf("invalid credentials format")
	}

	return credentials{
		username: creds[0],
		password: creds[1],
	}, nil
}
