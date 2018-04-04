package gomw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldRecoverFromPanicWithLog(t *testing.T) {
	next := &nextHandler{run: func() { panic("Some Error") }}
	l := &clogger{}
	rh := Recovery(l)(next)
	r, _ := http.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()

	rh.ServeHTTP(w, r)

	assert.True(t, next.called)
	assert.Equal(t, "panic: Some Error", l.log)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestShouldRecoverFromPanic(t *testing.T) {
	next := &nextHandler{run: func() { panic("Some Error") }}
	rh := Recovery(nil)(next)
	r, _ := http.NewRequest("GET", "/url", nil)
	w := httptest.NewRecorder()

	rh.ServeHTTP(w, r)

	assert.True(t, next.called)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type clogger struct {
	log string
}

func (l *clogger) Println(args ...interface{}) {
	l.log = fmt.Sprintf(args[0].(string), args[1:]...)
}
