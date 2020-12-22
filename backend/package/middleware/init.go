package middleware

import (
	"log"
	"net/http"
)

// Logger router
type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sw := statusRecorder{w, http.StatusOK}
	l.handler.ServeHTTP(&sw, r)
	log.Printf("%s %s %s", r.Method, r.URL.Path, http.StatusText(sw.status))
}

// NewLogger initiate logging router
func NewLogger(handler http.Handler) *Logger {
	return &Logger{
		handler: handler,
	}
}

type statusRecorder struct {
	w      http.ResponseWriter
	status int
}

func (s *statusRecorder) Header() http.Header {
	return s.w.Header()
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	return s.w.Write(b)
}

func (s *statusRecorder) WriteHeader(statusCode int) {
	s.status = statusCode
	s.w.WriteHeader(statusCode)
	return
}
