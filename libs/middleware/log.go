package middleware

import (
	"git.andresbott.com/Golang/carbon/libs/log"
	"net/http"
	"time"
)

// LoggingMiddleware logs every request using the provided logger
func LoggingMiddleware(next http.Handler, l log.LeveledStructuredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()
		l.DebugW("Request",
			"method", r.Method,
			"url", r.RequestURI,
			"time-code", timeStart,
		)

		logRespWriter := NewResponseWriter(w)
		next.ServeHTTP(logRespWriter, r)

		timeEnd := time.Now()
		timeDiff := timeEnd.Sub(timeStart)
		_ = timeDiff
		l.DebugW("Response: ",
			"duration", timeDiff,
			"url", r.RequestURI,
			"status-code", logRespWriter.statusCode,
			"time-code", timeStart,
		)
		return
	})
}
