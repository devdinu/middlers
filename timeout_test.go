package gomw

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type nextHandler struct {
	run        func()
	called     bool
	success    bool
	totalCalls int
	statusCode int
	Mutex      sync.Locker
}

func newMockNextHandler() *nextHandler {
	return &nextHandler{Mutex: &sync.Mutex{}}
}

func (h *nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.called = true
	if h.run != nil {
		h.run()
	}
	h.success = true
	h.totalCalls++
	if h.statusCode != 0 {
		w.WriteHeader(h.statusCode)
	}
}

func (h *nextHandler) Called() bool {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.called
}

func (h *nextHandler) Success() bool {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.success
}
func TestTimeoutMiddlewareCtxDeadlineExcceeded(t *testing.T) {
	next := &nextHandler{run: func() { time.Sleep(10 * time.Millisecond) }, Mutex: new(sync.Mutex)}
	d := 1 * time.Millisecond
	handler := Timeout(d)(next)
	req, _ := http.NewRequest("GET", "/some/url", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, next.Called())
	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

func TestTimeoutMiddlewareHandlerSucceedsBeforeDeadline(t *testing.T) {
	next := &nextHandler{run: func() { time.Sleep(10 * time.Millisecond) }, Mutex: &sync.Mutex{}}
	d := 15 * time.Millisecond
	handler := Timeout(d)(next)
	req, _ := http.NewRequest("GET", "/some/url/success", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, next.Called())
	assert.True(t, next.Success())
	assert.Equal(t, http.StatusOK, w.Code)
}
