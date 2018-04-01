package gomw

import (
	"net/http"
	"strings"
)

type reporter interface {
	Increment(string)
}

type Middleware func(http.Handler) http.Handler

func StatsReporter(rep reporter) Middleware {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			rw := NewResponseWriter(w)
			next.ServeHTTP(rw, r)
			go func() {
				status := strings.Replace(strings.ToLower(http.StatusText(rw.StatusCode)), " ", "_", -1)
				url := strings.TrimLeft(strings.Replace(r.URL.Path, "/", "_", -1), "_")
				rep.Increment(url + "_" + status)
			}()
		}
		return http.HandlerFunc(h)
	}
}
