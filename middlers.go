package gomw

import "net/http"

type Config struct {
	predicate func(*http.Request) bool
	rlogger   logger
}

type Option func(*Config)

func New(next http.Handler, opts ...Option) http.Handler {
	cfg := &Config{}
	handler := next
	for _, o := range opts {
		o(cfg)
	}
	if cfg.predicate != nil {
		handler = Filter(cfg.predicate)(handler)
	}
	if cfg.rlogger != nil {
		handler = RequestLogger(handler, cfg.rlogger)
	}
	return handler
}

func Predicate(f func(*http.Request) bool) Option {
	return func(c *Config) {
		c.predicate = f
	}
}

func Logger(l logger) Option {
	return func(c *Config) {
		c.rlogger = l
	}
}
