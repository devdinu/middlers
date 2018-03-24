package gomw

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareToCallNext(t *testing.T) {
	next := &nextHandler{}
	testcases := []struct {
		description string
		handler     http.Handler
	}{
		{
			description: "call next from TimeoutMiddleware",
			handler:     Timeout(next, time.Second),
		},
		{
			description: "call next from Filter Middleware",
			handler:     Filter(func(r *http.Request) bool { return true }, next),
		},
		{
			description: "call next from logger middleware",
			handler:     RequestLogger(next),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "/some/url", nil)
			w := httptest.NewRecorder()
			tc.handler.ServeHTTP(w, r)
		})
	}
	assert.Equal(t, len(testcases), next.totalCalls)
}
