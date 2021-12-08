package log

import (
	"net/http"
	"time"
)

// LoggingMiddleware logs every request using the provided logger
func LoggingMiddleware(next http.Handler, l LeveledStructuredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()

		logRespWriter := NewResponseWriter(w)
		next.ServeHTTP(logRespWriter, r)

		timeEnd := time.Now()
		timeDiff := timeEnd.Sub(timeStart)

		if IsStatusError(logRespWriter.statusCode) {
			l.Error("",
				"method", r.Method,
				"url", r.RequestURI,
				"duration", timeDiff,
				"response-code", logRespWriter.statusCode,
				"user-agent", r.UserAgent(),
				"referer", r.Referer(),
				"ip", ReadUserIP(r),
				"req-id", r.Header.Get("Request-Id"),
			)
		}
		l.Debug("",
			"method", r.Method,
			"url", r.RequestURI,
			"duration", timeDiff,
			"response-code", logRespWriter.statusCode,
			"user-agent", r.UserAgent(),
			"referer", r.Referer(),
			"ip", ReadUserIP(r),
			"req-id", r.Header.Get("Request-Id"),
		)
		return
	})
}

// https://stackoverflow.com/questions/27234861/correct-way-of-getting-clients-ip-addresses-from-http-request
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
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

func IsStatusError(statusCode int) bool {
	return statusCode < 600 && statusCode >= 500
}
