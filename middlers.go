package gomw

import (
	"net/http"
	"time"
)

type config struct {
	predicate func(*http.Request) bool
	rlogger   logger
	Store
	RateLimitConfig
	timeout time.Duration
}

type Option func(*config)

func New(next http.Handler, opts ...Option) http.Handler {
	cfg := &config{}
	handler := next
	for _, o := range opts {
		o(cfg)
	}
	if cfg.predicate != nil {
		handler = Filter(cfg.predicate)(handler)
	}
	if cfg.timeout != 0 {
		handler = Timeout(cfg.timeout)(handler)
	}
	if cfg.Store != nil {
		handler = Ratelimit(cfg.Store, cfg.RateLimitConfig)(handler)
	}
	if cfg.rlogger != nil {
		handler = RequestLogger(handler, cfg.rlogger)
	}
	return handler
}

func Predicate(f func(*http.Request) bool) Option {
	return func(c *config) {
		c.predicate = f
	}
}

func Logger(l logger) Option {
	return func(c *config) {
		c.rlogger = l
	}
}

func InMemoryRateLimit(cfg RateLimitConfig) Option {
	return func(c *config) {
		c.RateLimitConfig = cfg
		c.Store = InMemoryStore()
	}
}

func RateLimitter(s Store, cfg RateLimitConfig) Option {
	return func(c *config) {
		c.Store = s
		c.RateLimitConfig = cfg
		if s == nil {
			c.Store = InMemoryStore()
		}
	}
}

func Timed(d time.Duration) Option {
	return func(c *config) {
		c.timeout = d
	}
}
