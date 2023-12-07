package ratelimit

import (
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
	if rate <= 0 {
		rate = DefaultRate
	}
	if window <= 0 {
		window = DefaultWindow
	}

	bs := &BucketStore{
		buckets:   map[string]chan token{},
		bucketLen: rate,
		interval:  window / time.Duration(rate),
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
