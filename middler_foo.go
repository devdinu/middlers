package gomw

import (
	"fmt"
	"net/http"
	"time"
)

//Filter, Logging, Stats, Timeout

type COption func(*MwConfig)

type CC struct {
	Timeout      time.Duration
	StatsEnabled bool
}

type MwConfig struct {
	loggingEnabled bool
	timeout        time.Duration
	statsAddress   string
	// customerlogger
	dbTimeout int
}

func MwWithTimeout(t time.Duration) COption {
	return func(c *MwConfig) {
		fmt.Println("setting timeout")
		c.timeout = t
	}
}

func MwWithStats(st string) COption {
	return func(c *MwConfig) {
		fmt.Println("setting statsd")
		if st != "" {
			c.statsAddress = st
			return
		}
		c.statsAddress = "localhost"
	}
}

func NewMiddleware(next http.Handler, opts ...COption) http.Handler {
	cc := &MwConfig{}
	for _, opt := range opts {
		opt(cc)
	}
	fmt.Println("actual config:", cc)
	return http.NotFoundHandler()
}

//NewConfig() MWConfig
//NewMiddleware(next, MwConfig{})

//NewMiddleware(next, nil)
