package gomw_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/devdinu/middlers"
	"github.com/stretchr/testify/assert"
)

type logger struct {
	called bool
}

func (l *logger) Println(args ...interface{}) {
	l.called = true
}

type rateLimitStore struct {
	called bool
}

func (r *rateLimitStore) Get(s string) int { r.called = true; return 0 }
func (r *rateLimitStore) Incr(s string)    {}
func (r *rateLimitStore) Reset(s string)   {}

func TestMiddler(t *testing.T) {
	var called, predicateCalled bool
	predicate := func(r *http.Request) bool {
		predicateCalled = true
		return true
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	lg, rstore := &logger{}, &rateLimitStore{}
	rcfg := gomw.RateLimitConfig{
		MaxRequests:     1,
		RequestKey:      func(r *http.Request) string { return "some" },
		TimeWindowReset: 1000,
	}

	middler := gomw.New(handler,
		gomw.Predicate(predicate),
		gomw.Logger(lg),
		gomw.RateLimitter(rstore, rcfg),
		gomw.Timed(2000),
	)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/url", nil)

	middler.ServeHTTP(w, r)

	assert.True(t, predicateCalled, "Should have called predicate")
	assert.True(t, called, "Should call actual handler")
	assert.True(t, lg.called, "logger should've been called")
	assert.True(t, rstore.called, "ratelimitter should've been called")
	assert.Equal(t, w.Code, 200)
}

func TestBla(t *testing.T) {
	cfg := gomw.CC{Timeout: time.Duration(100)}

	m := gomw.NewMiddleware(http.NotFoundHandler(),
		gomw.MwWithTimeout(time.Duration(10)),
		gomw.MwWithStats("127.0.0.1"),
	)

	fmt.Println(cfg, m)
	//	TimeoutF := func(c *MwConfig) {
	//		c.timeout = time.Duration(10)
	//	}
}
