package gomw

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Store interface {
	Get(string) int
	Incr(string)
	Reset(string)
}

type RateLimitConfig struct {
	MaxRequests     int
	TimeWindowReset time.Duration
	RequestKey      func(*http.Request) string
}

const rateLimitResetHeader = "X-Ratelimit-Reset"

func Ratelimit(s Store, cfg RateLimitConfig) Middleware {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			if ratelimit(s, cfg, r) {
				w.Header().Set(rateLimitResetHeader, fmt.Sprintf("%.2f", cfg.TimeWindowReset.Seconds()))
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
}

func ratelimit(s Store, cfg RateLimitConfig, r *http.Request) bool {
	if cfg.RequestKey == nil || s == nil || cfg.MaxRequests == 0 || cfg.TimeWindowReset == 0 {
		return false
	}
	key := cfg.RequestKey(r)
	defer s.Incr(key)
	val := s.Get(key)
	if val == 0 {
		time.AfterFunc(cfg.TimeWindowReset, func() { s.Reset(key) })
	}
	if val < cfg.MaxRequests {
		return false
	}
	return true
}

type mapStore struct {
	sync.Mutex
	record map[string]int
}

func (ms *mapStore) Get(k string) int {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	return ms.record[k]
}

func (ms *mapStore) Reset(k string) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	delete(ms.record, k)
}

func (ms *mapStore) Incr(k string) {
	ms.Mutex.Lock()
	defer ms.Mutex.Unlock()
	ms.record[k] += 1
}

func InMemoryStore() Store {
	return new(mapStore)
}
