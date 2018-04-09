package gomw

import "net/http"

func Filter(predicate func(*http.Request) bool) Middleware {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			if predicate(r) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
		return http.HandlerFunc(mw)
	}
}
