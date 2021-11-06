package middleware

import (
	"git.andresbott.com/Golang/carbon/libs/log"
	"net/http"
	"strconv"
	"time"
)

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

// LogResponseWriter allows to o get access to the response within the middleware function
type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	// WriteHeader(int) is not called if the response is 200 (implicit response code) so it needs to be the default
	return &LogResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}
func (r *LogResponseWriter) StatusCode() int {
	return r.statusCode
}

func (r *LogResponseWriter) StatusCodeStr() string {
	return strconv.Itoa(r.statusCode)
}

// Write returns underlying Write result, while counting data size
func (r *LogResponseWriter) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r *LogResponseWriter) WriteHeader(code int) {
	r.statusCode = code
	// avoid superfluous status code warning
	if code != http.StatusOK {
		r.ResponseWriter.WriteHeader(code)
	}

}
