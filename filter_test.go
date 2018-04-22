package gomw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterMiddleware(t *testing.T) {
	block := func(r *http.Request) bool { return false }
	allow := func(r *http.Request) bool { return true }
	r, _ := http.NewRequest("GET", "/some/url", nil)

	t.Run("should not call next on false predicate", func(t *testing.T) {
		next := newMockNextHandler()
		h := Filter(block)(next)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		assert.False(t, next.called)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should call next on true predicate", func(t *testing.T) {
		next := newMockNextHandler()
		h := Filter(allow)(next)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		assert.True(t, next.called)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
