package gomw

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, StatusCode: 200}
}
