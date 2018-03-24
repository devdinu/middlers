package gomw

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type nextHandler struct {
	run     func()
	called  bool
	success bool
}

func (h *nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	h.run()
	h.success = true
}

func TestTimeoutMiddlewareCtxDeadlineExcceeded(t *testing.T) {
	next := &nextHandler{run: func() { time.Sleep(10 * time.Millisecond) }}
	d := 5 * time.Millisecond
	handler := Timeout(next, d)
	req, _ := http.NewRequest("GET", "/some/url", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, next.called)
	assert.False(t, next.success)
	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

func TestTimeoutMiddlewareHandlerSucceedsBeforeDeadline(t *testing.T) {
	next := &nextHandler{run: func() { time.Sleep(10 * time.Millisecond) }}
	d := 15 * time.Millisecond
	handler := Timeout(next, d)
	req, _ := http.NewRequest("GET", "/some/url/success", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, next.called)
	assert.True(t, next.success)
	assert.Equal(t, http.StatusOK, w.Code)
}
