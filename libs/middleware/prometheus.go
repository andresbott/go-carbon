package middleware

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"time"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests for different paths, methods, status codes",
	},
		[]string{
			"type",
			"status",
			"method",
			"addr",
			"isError",
			//"errorMessage",
		},
	)
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		respWriter := NewResponseWriter(w)

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timeStart := time.Now()
		next.ServeHTTP(respWriter, r)

		statusCodeStr := respWriter.StatusCodeStr()
		isErrorStr := strconv.FormatBool(IsStatusError(respWriter.statusCode))

		httpDuration.With(prometheus.Labels{
			"type":    r.Proto,
			"status":  statusCodeStr,
			"method":  r.Method,
			"addr":    path,
			"isError": isErrorStr,
		}).Observe(time.Since(timeStart).Seconds())
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
func (r *ResponseWriter) StatusCode() int {
	return r.statusCode
}

func (r *ResponseWriter) StatusCodeStr() string {
	return strconv.Itoa(r.statusCode)
}

// Write returns underlying Write result, while counting data size
func (r *ResponseWriter) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	return n, err
}

func (r *ResponseWriter) WriteHeader(code int) {
	r.statusCode = code
	// avoid superflous status code warning
	if code != http.StatusOK {
		r.ResponseWriter.WriteHeader(code)
	}

}

func IsStatusError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 400
}
