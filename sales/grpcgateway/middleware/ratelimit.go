package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

const (
	requestsPerMinute = 30
	intervalInSeconds = 60 // 60 seconds in a minute
)

// Calculate the rate as requests per second
var ratePerSecond = rate.Limit(requestsPerMinute / intervalInSeconds)

// The burst size is the maximum number of requests allowed in a single burst, set to 5
var burstSize = 5

// Create a rate limiter instance (adjust the parameters as per your requirements)
var limiter = rate.NewLimiter(ratePerSecond, burstSize) // 1 request per second with a burst of 5

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
