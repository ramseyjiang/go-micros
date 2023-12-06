package ratelimit

import "net/http"

const defaultKey = "test rate limit" // The key can be remoteIP, userID or something else.

func Impl(bucketStore *BucketStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, _, _, _ := bucketStore.Take(defaultKey)
			if !allowed {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
