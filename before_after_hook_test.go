package gomw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeforeCall(t *testing.T) {
	var beforeCalled, afterCalled bool
	next := &nextHandler{}
	before := func() { beforeCalled = true }
	after := func() { afterCalled = true }

	r, _ := http.NewRequest("GET", "/some/url", nil)
	w := httptest.NewRecorder()

	ExecutionHooks(before, after)(next).ServeHTTP(w, r)

	assert.True(t, next.called)
	assert.True(t, beforeCalled)
	assert.True(t, afterCalled)
	assert.Equal(t, 200, w.Code)
}

func TestShouldCallAfterEvenIfHandlerPanics(t *testing.T) {
	var beforeCalled, afterCalled bool
	next := &nextHandler{run: func() { panic("panic error") }}
	before := func() { beforeCalled = true }
	after := func() { afterCalled = true }

	r, _ := http.NewRequest("GET", "/some/url", nil)
	w := httptest.NewRecorder()

	assert.Panics(t,
		func() {
			ExecutionHooks(before, after)(next).ServeHTTP(w, r)
		},
	)

	assert.True(t, next.called)
	assert.True(t, beforeCalled, "before should've been called")
	assert.True(t, afterCalled, "after should've been called")
	assert.Equal(t, 200, w.Code)
}

func TestShouldCallBeforeAfterInOrder(t *testing.T) {
	var beforeCalled, afterCalled bool
	var order []string
	next := &nextHandler{run: func() { order = append(order, "handler") }}
	before := func() {
		order = append(order, "before")
		beforeCalled = true
	}
	after := func() {
		order = append(order, "after")
		afterCalled = true
	}

	r, _ := http.NewRequest("GET", "/some/url", nil)
	w := httptest.NewRecorder()

	ExecutionHooks(before, after)(next).ServeHTTP(w, r)

	assert.True(t, next.called)
	assert.True(t, beforeCalled, "before should've been called")
	assert.True(t, afterCalled, "after should've been called")
	assert.Equal(t, []string{"before", "handler", "after"}, order)
	assert.Equal(t, 200, w.Code)
}
