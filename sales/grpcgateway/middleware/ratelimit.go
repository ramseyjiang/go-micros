package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

const (
	interval  = 1
	frequency = 5
)

// Create a rate limiter instance (adjust the parameters as per your requirements)
var limiter = rate.NewLimiter(interval, frequency) // 1 request per second with a burst of 5

// RateLimit enforces rate limiting on incoming requests using build-in library.
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
