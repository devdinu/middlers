package gomw

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReporter struct {
	mock.Mock
}

func (m *MockReporter) Increment(s string) {
	m.Called(s)
}

func TestShouldIncrementUrlStats(t *testing.T) {
	next := &nextHandler{statusCode: http.StatusUnauthorized, Mutex: &sync.Mutex{}}
	rmock := new(MockReporter)
	done := make(chan bool)
	rmock.On("Increment", "some_path_unauthorized").Run(func(mock.Arguments) { done <- true })
	r, _ := http.NewRequest("GET", "/some/path", nil)

	mw := StatsReporter(rmock)(next)
	mw.ServeHTTP(httptest.NewRecorder(), r)

	<-done
	assert.True(t, next.called)
	rmock.AssertExpectations(t)
}
