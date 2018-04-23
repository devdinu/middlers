package gomw

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimitPassRequest(t *testing.T) {
	cfg := RateLimitConfig{MaxRequests: 2, TimeWindowReset: 1 * time.Second}
	next := newMockNextHandler()
	rl := Ratelimit(InMemoryStore(), cfg)(next)
	calls := 2

	for i := 0; i < calls; i++ {
		r, _ := http.NewRequest("GET", "/url", nil)
		w := httptest.NewRecorder()
		rl.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.Empty(t, w.Header().Get("X-RateLimit-Reset"))
	}

	assert.Equal(t, calls, next.totalCalls)
	assert.True(t, next.Success())
	assert.True(t, next.Called())
}

func TestRateLimitPassWhenConfigIsNotProper(t *testing.T) {
	ms := new(mockStore)
	next := newMockNextHandler()
	ms.On("Get", mock.AnythingOfType("string")).Return(5)
	ms.On("Incr", mock.AnythingOfType("string"))
	ms.On("Reset", mock.AnythingOfType("string"))

	testcases := []struct {
		description string
		cfg         RateLimitConfig
		store       Store
	}{
		{
			"store is empty",
			RateLimitConfig{MaxRequests: 1, TimeWindowReset: 1, RequestKey: func(*http.Request) string { return "" }},
			Store(nil),
		},
		{
			"reset timer config is empty",
			RateLimitConfig{MaxRequests: 1, RequestKey: func(*http.Request) string { return "" }},
			ms,
		},
		{
			"function fetch key from request config is empty",
			RateLimitConfig{MaxRequests: 1, TimeWindowReset: 2},
			ms,
		},
		{
			"RateLimit max requests config is 0",
			RateLimitConfig{TimeWindowReset: 1, RequestKey: func(*http.Request) string { return "" }},
			ms,
		},
	}

	for _, tc := range testcases {
		rl := Ratelimit(tc.store, tc.cfg)(next)
		r, _ := http.NewRequest("GET", "/url", nil)
		w := httptest.NewRecorder()
		rl.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code, tc.description)
		assert.Empty(t, w.Header().Get("X-RateLimit-Reset"), tc.description)
	}

	assert.Equal(t, len(testcases), next.totalCalls)
	assert.True(t, next.Success())
	assert.True(t, next.Called())

}

func TestRateLimitFailure(t *testing.T) {
	cfg := RateLimitConfig{MaxRequests: 1, TimeWindowReset: 2000 * time.Millisecond, RequestKey: func(r *http.Request) string { return "request_key" }}
	next, ms := newMockNextHandler(), new(mockStore)
	rl := Ratelimit(ms, cfg)(next)
	ms.On("Get", "request_key").Return(2).Once()
	ms.On("Incr", "request_key").Once()
	r, _ := http.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()

	rl.ServeHTTP(w, r)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "2.00", w.Header().Get("X-RateLimit-Reset"))
	assert.Empty(t, next.totalCalls)
	assert.False(t, next.Success())
	assert.False(t, next.Called())
	ms.AssertExpectations(t)
}

func TestRateLimitAfterExpiry(t *testing.T) {
	cfg := RateLimitConfig{MaxRequests: 1, TimeWindowReset: 50 * time.Millisecond, RequestKey: func(r *http.Request) string { return r.URL.Path }}
	next, s := newMockNextHandler(), &mapStore{Mutex: sync.Mutex{}, record: make(map[string]int)}
	rl := Ratelimit(s, cfg)(next)
	r, _ := http.NewRequest("GET", "/url", nil)

	w := httptest.NewRecorder()

	rl.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Empty(t, w.Header().Get("X-RateLimit-Reset"))
	assert.Equal(t, 1, s.Get(r.URL.Path), "Store should have the record of request")

	w = httptest.NewRecorder()
	rl.ServeHTTP(w, r)
	assert.Equal(t, 429, w.Code, "Should fail the second immediate request")
	assert.Equal(t, "0.05", w.Header().Get("X-RateLimit-Reset"))
	assert.Equal(t, 2, s.Get(r.URL.Path), "Store should have the record of request")

	time.Sleep(70 * time.Millisecond)

	w = httptest.NewRecorder()
	rl.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code, "Should succeed after the timeout")
	assert.Empty(t, w.Header().Get("X-RateLimit-Reset"))
	assert.Equal(t, 1, s.Get(r.URL.Path), "Store should have the record of request")

	assert.Equal(t, 2, next.totalCalls)
	assert.True(t, next.Success())
	assert.True(t, next.Called())
}

type mockStore struct{ mock.Mock }

func (m *mockStore) Reset(s string)   { m.Called(s) }
func (m *mockStore) Incr(s string)    { m.Called(s) }
func (m *mockStore) Get(s string) int { return m.Called(s).Int(0) }
