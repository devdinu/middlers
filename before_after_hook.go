package gomw

import (
	"net/http"
)

func ExecutionHooks(before, after func()) Middleware {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			if before != nil {
				before()
			}
			if after != nil {
				defer after()
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
}
