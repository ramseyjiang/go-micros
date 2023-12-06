package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	// Define a very restrictive rate limit for testing
	testRate := 2
	testWindow := 1 * time.Minute

	// Initialize the test BucketStore with the test rate and window
	testBucketStore := NewBucketStore(testRate, testWindow)

	// Dummy handler for testing
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the dummy handler with the RateLimit middleware using the test bucket store
	rateLimitedHandler := Impl(testBucketStore)(dummyHandler)

	tests := []struct {
		name            string
		numRequests     int
		requestInterval time.Duration
		wantStatusCodes []int
	}{
		{
			name:            "WithinRateLimit",
			numRequests:     1,
			requestInterval: 500 * time.Millisecond, // Half a second between requests
			wantStatusCodes: []int{http.StatusOK, http.StatusOK},
		},
		{
			name:            "ExceedRateLimit",
			numRequests:     2,                      // Two requests to test the limit
			requestInterval: 500 * time.Millisecond, // Half a second between requests
			wantStatusCodes: []int{http.StatusOK, http.StatusTooManyRequests},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.numRequests; i++ {
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()
				rateLimitedHandler.ServeHTTP(rr, req)

				if i < len(tt.wantStatusCodes) {
					if status := rr.Code; status != tt.wantStatusCodes[i] {
						t.Errorf("handler returned wrong status code for request %d: got %v want %v", i+1, status, tt.wantStatusCodes[i])
					}
				}

				time.Sleep(tt.requestInterval)
			}
		})
	}
}
