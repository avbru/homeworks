package internalhttp

import (
	"net/http"
	"time"

	"github.com/avbru/homeworks/hw12_13_14_15_calendar/internal/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := statusRecorder{w, 200}
		next.ServeHTTP(&rec, r)
		duration := time.Since(startTime)
		logger.Logger.Info().Msgf("%s %s %s %s %s %d %s %s",
			r.RemoteAddr,
			startTime.String(),
			r.Method,
			r.URL.Path,
			r.Proto,
			rec.status,
			duration.String(),
			r.UserAgent())
	})
}
