package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"krizanauskas.github.com/mvp-proxy/internal/services"
)

type BandwidthLimitMiddleware struct {
	bandwidthControler services.UserBandwidthControllerI
}

func NewBandwidthLimitMiddleware(bandwidthController services.UserBandwidthControllerI) BandwidthLimitMiddleware {
	return BandwidthLimitMiddleware{
		bandwidthController,
	}
}

func (m BandwidthLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := r.Context().Value(userContextKey).(*User)
		if !ok {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if m.bandwidthControler.GetAvailableBandwidth(user.Username) <= 0 {
			if r.Method == http.MethodConnect {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, "Sorry, you have exceeded your bandwidth limit..")
				return
			}

			w.WriteHeader(http.StatusForbidden)
			file, err := os.Open("web/templates/bandwidth_exceeded.html")
			if err != nil {
				return
			}

			defer file.Close()

			_, _ = io.Copy(w, file)

			return
		}

		next.ServeHTTP(w, r)
	})
}
