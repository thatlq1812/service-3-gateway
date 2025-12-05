package middleware

import (
	"context"
	"net/http"
	"time"
)

// TimeoutMiddleware adds request timeout to prevent hanging requests
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Pass context with timeout to next handler
			r = r.WithContext(ctx)

			// Channel to signal handler completion
			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
				// Request completed successfully
				return
			case <-ctx.Done():
				// Timeout occurred
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte(`{"code":"504","message":"request timeout: service took too long to respond"}`))
			}
		})
	}
}
