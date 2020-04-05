package gomw

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type logger interface {
	Println(...interface{})
}

type requestLog struct {
	start  time.Time
	end    time.Time
	r      *http.Request
	status int
}

func (rl *requestLog) String() string {
	//TODO: could log status if its negroni.ResponseWriter
	return fmt.Sprintf(`{"method: %s, url: %s, status: %d, requested_at: %v, response_at: %v, duration_ms: %v}`,
		rl.r.Method, rl.r.URL.Path, rl.status, rl.start.Format(time.RFC3339), rl.end.Format(time.RFC3339), rl.end.Sub(rl.start))
}

func RequestLogger(next http.Handler, optloggers ...logger) http.HandlerFunc {
	var clogger logger
	if len(optloggers) > 0 {
		clogger = optloggers[0]
	} else {
		clogger = log.New(os.Stdout, "", 0)
	}
	mw := func(w http.ResponseWriter, r *http.Request) {
		rl := &requestLog{}
		rl.start = time.Now()
		rw := NewResponseWriter(w)
		time.Sleep(time.Second)

		next.ServeHTTP(rw, r)

		rl.end = time.Now()
		rl.r = r
		rl.status = rw.StatusCode
		clogger.Println(rl.String())
	}
	return http.HandlerFunc(mw)
}
