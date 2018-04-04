package gomw

import "net/http"

func Recovery(l logger) Middleware {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if l != nil {
						l.Println("panic: %v", err)
					}
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(h)
	}
}
