package middlewares_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"krizanauskas.github.com/mvp-proxy/internal/middlewares"
	mock_services "krizanauskas.github.com/mvp-proxy/tests/mocks"
)

type User struct {
	Username string
}

const validUsername = "testuser"
const validPassword = "password"

func TestBasicAuthMiddleware_SuccessfulAuthentication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockAuthServiceI(ctrl)

	middleware := middlewares.NewBasicAuthMiddleware(mockAuthService)

	credentials := base64.StdEncoding.EncodeToString([]byte(validUsername + ":" + validPassword))

	mockAuthService.EXPECT().Authenticate(validUsername, validPassword).Return(true)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Proxy-Authorization", "Basic "+credentials)

	recorder := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(middlewares.UserContextKey).(*middlewares.User)
		assert.Equal(t, validUsername, user.Username)
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestBasicAuthMiddleware_MissingCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockAuthServiceI(ctrl)
	middleware := middlewares.NewBasicAuthMiddleware(mockAuthService)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusProxyAuthRequired, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "failed to parse credentials")
}

func TestBasicAuthMiddleware_InvalidCredentialsFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockAuthServiceI(ctrl)

	middleware := middlewares.NewBasicAuthMiddleware(mockAuthService)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Proxy-Authorization", "Bearer xyz")

	recorder := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusProxyAuthRequired, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "failed to parse credentials")
}

func TestBasicAuthMiddleware_FailedAuthentication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockAuthServiceI(ctrl)

	middleware := middlewares.NewBasicAuthMiddleware(mockAuthService)

	credentials := base64.StdEncoding.EncodeToString([]byte("bad_user:bad_pass"))
	mockAuthService.EXPECT().Authenticate("bad_user", "bad_pass").Return(false)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Proxy-Authorization", "Basic "+credentials)

	recorder := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusProxyAuthRequired, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "failed to authenticate")
}
