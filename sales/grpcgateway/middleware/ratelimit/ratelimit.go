package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

const (
	DefaultRate   = 5 // Default requests per minute
	DefaultWindow = 1 * time.Minute
)

type token struct{}

// BucketStore describes a Token Bucket store
type BucketStore struct {
	sync.Mutex // guards buckets
	buckets    map[string]chan token
	bucketLen  int
	interval   time.Duration
	Reset      time.Time
}

// NewBucketStore creates new in-memory token bucket store.
// The purpose of this store is to allow for local rate limiting backoff,
// before the datastore is hammered during DDoS attacks.
func NewBucketStore(rate int, window time.Duration) *BucketStore {
	bs := &BucketStore{
		buckets:   map[string]chan token{},
		bucketLen: rate,
		interval:  window,
	}
	bs.startTicker()
	return bs
}

func (s *BucketStore) startTicker() {
	tick := time.NewTicker(s.interval)
	go func() {
		for t := range tick.C {
			s.Lock()
			s.Reset = t.Add(s.interval)
			for key, bucket := range s.buckets {
				select {
				case <-bucket:
				default:
					delete(s.buckets, key)
				}
			}
			s.Unlock()
		}
	}()
}

// InitRate initialises the rate of the Token Bucket store
func (s *BucketStore) InitRate(rate int, window time.Duration) {
	if rate == 0 {
		rate = DefaultRate
	}
	if window.Nanoseconds() == 0 {
		window = DefaultRate * time.Minute
	}
	s.bucketLen = rate
	s.Reset = time.Now()
	s.interval = time.Duration(int(window) / rate)

	go func() {
		interval := time.Duration(int(window) / rate)
		tick := time.NewTicker(interval)
		for t := range tick.C {
			s.Lock()
			s.Reset = t.Add(interval)
			for key, bucket := range s.buckets {
				select {
				case <-bucket:
				default:
					delete(s.buckets, key)
				}
			}
			s.Unlock()
		}
	}()
}

// GetRate gets the Rate (the bucketLen)
func (s *BucketStore) GetRate() string {
	return fmt.Sprintf("%v", s.bucketLen)
}

// GetInterval gets the interval
func (s *BucketStore) GetInterval() string {
	return fmt.Sprintf("%v", s.interval)
}

// Take implements TokenBucketStore interface.
// It takes token from a bucket referenced by a given key, if available.
func (s *BucketStore) Take(key string) (bool, int, time.Time, error) {
	s.Lock()
	bucket, ok := s.buckets[key]
	if !ok {
		bucket = make(chan token, s.bucketLen)
		s.buckets[key] = bucket
	}
	s.Unlock()
	select {
	case bucket <- token{}:
		return true, cap(bucket) - len(bucket), s.Reset, nil
	default:
		return false, 0, s.Reset, nil
	}
}
