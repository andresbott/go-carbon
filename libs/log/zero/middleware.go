package zero

import (
	"git.andresbott.com/Golang/carbon/libs/http/extra/respstatus"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

// LoggingMiddleware logs every request using the provided logger
func LoggingMiddleware(next http.Handler, l *zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()

		respWriter := respstatus.NewWriter(w)
		next.ServeHTTP(respWriter, r)

		timeEnd := time.Now()
		timeDiff := timeEnd.Sub(timeStart)

		if respstatus.IsStatusError(respWriter.StatusCode()) {
			l.Error().
				Str("method", r.Method).
				Str("url", r.RequestURI).
				Dur("durations", timeDiff).
				Int("response-code", respWriter.StatusCode()).
				Str("user-agent", r.UserAgent()).
				Str("referer", r.Referer()).
				Str("ip", ReadUserIP(r)).
				Str("req-id", r.Header.Get("Request-Id")).
				Msg("")
		} else {
			l.Info().
				Str("method", r.Method).
				Str("url", r.RequestURI).
				Dur("durations", timeDiff).
				Int("response-code", respWriter.StatusCode()).
				Str("user-agent", r.UserAgent()).
				Str("referer", r.Referer()).
				Str("ip", ReadUserIP(r)).
				Str("req-id", r.Header.Get("Request-Id")).
				Msg("")
		}
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
