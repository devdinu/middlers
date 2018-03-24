package gomw

import (
	"context"
	"net/http"
	"time"
)

func Timeout(next http.Handler, timeout time.Duration) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		success := make(chan bool)
		defer cancel()
		go func() {
			next.ServeHTTP(w, r.WithContext(ctx))
			success <- true
		}()
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusGatewayTimeout)
		case <-success:
		}
	}
	return http.HandlerFunc(mw)
}
