package log

import (
	"net/http"
	"time"
)

// LoggingMiddleware logs every request using the provided logger
func LoggingMiddleware(next http.Handler, l LeveledStructuredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()
		l.DebugW("Request",
			"method", r.Method,
			"url", r.RequestURI,
		)

		logRespWriter := NewResponseWriter(w)
		next.ServeHTTP(logRespWriter, r)

		timeEnd := time.Now()
		timeDiff := timeEnd.Sub(timeStart)
		l.DebugW("Response: ",
			"duration", timeDiff,
			"url", r.RequestURI,
			"status-code", logRespWriter.statusCode,
		)

		l.InfoW("Request: ",
			"method", r.Method,
			"url", r.RequestURI,
			"duration", timeDiff,
			"response-code", logRespWriter.statusCode,
		)
		return
	})
}

// ResponseWriter allows to get the status code of the response in the middleware
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	// WriteHeader(int) is not called if the response is 200 (implicit response code) so it needs to be the default
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}
