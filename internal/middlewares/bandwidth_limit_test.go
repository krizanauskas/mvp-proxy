package middlewares_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"krizanauskas.github.com/mvp-proxy/internal/middlewares"
	mock_services "krizanauskas.github.com/mvp-proxy/tests/mocks"
)

func TestBandwidthLimitMiddleware_UserExceedsLimitMethodConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBandwidthController := mock_services.NewMockUserBandwidthControllerI(ctrl)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodConnect, "/", nil)

	user := &middlewares.User{Username: "testUser"}
	ctx := context.WithValue(req.Context(), middlewares.UserContextKey, user)
	req = req.WithContext(ctx)

	mockBandwidthController.EXPECT().GetAvailableBandwidth(user.Username).Return(0)

	middleware := middlewares.NewBandwidthLimitMiddleware(mockBandwidthController)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Sorry, you have exceeded your bandwidth limit")
}

func TestBandwidthLimitMiddleware_UserExceedsLimitMethodGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBandwidthController := mock_services.NewMockUserBandwidthControllerI(ctrl)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	user := &middlewares.User{Username: "testUser"}
	ctx := context.WithValue(req.Context(), middlewares.UserContextKey, user)
	req = req.WithContext(ctx)

	mockBandwidthController.EXPECT().GetAvailableBandwidth(user.Username).Return(0)

	middleware := middlewares.NewBandwidthLimitMiddleware(mockBandwidthController)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestBandwidthLimitMiddleware_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBandwidthController := mock_services.NewMockUserBandwidthControllerI(ctrl)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	middleware := middlewares.NewBandwidthLimitMiddleware(mockBandwidthController)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.Middleware(nextHandler).ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "User not found")
}
