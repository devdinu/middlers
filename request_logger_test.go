package gomw

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLoggerWithCustomLogger(t *testing.T) {
	l := log.New(os.Stdout, "customlogger: ", 0)
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
	h := RequestLogger(next, l)
	req, _ := http.NewRequest("GET", "/some/url", nil)

	h(httptest.NewRecorder(), req)
	assert.True(t, called)
}

func TestRequestLoggerWithNoLogger(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
	h := RequestLogger(next)
	req, _ := http.NewRequest("GET", "/some/url", nil)

	h(httptest.NewRecorder(), req)
	assert.True(t, called)
}
