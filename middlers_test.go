package gomw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devdinu/middlers"
	"github.com/stretchr/testify/assert"
)

type logger struct {
	called bool
}

func (l *logger) Println(args ...interface{}) {
	l.called = true
}

func TestMiddlerPredicate(t *testing.T) {
	var called, predicateCalled bool
	predicate := func(r *http.Request) bool {
		predicateCalled = true
		return true
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	lg := &logger{}

	middler := gomw.New(handler,
		gomw.Predicate(predicate),
		gomw.Logger(lg),
	)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/url", nil)

	middler.ServeHTTP(w, r)

	assert.True(t, predicateCalled, "Should have called predicate")
	assert.True(t, called, "Should call actual handler")
	assert.True(t, lg.called)
	assert.Equal(t, w.Code, 200)
}
