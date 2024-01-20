package respstatus

import (
	"net/http"
	"strconv"
)

// StatWriter is a wrapper to a httpResponse writer that allows to intercept and
// extract the status code that the upstream code has defined
type StatWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewWriter will return a pointer to a response writer
func NewWriter(w http.ResponseWriter) *StatWriter {
	// WriteHeader(int) is not called if the response is 200 (implicit response code) so it needs to be the default
	return &StatWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *StatWriter) StatusCode() int {
	return r.statusCode
}

func (r *StatWriter) StatusCodeStr() string {
	return strconv.Itoa(r.statusCode)
}

//// Write returns underlying Write result, while counting data size
//func (r *StatWriter) Write(b []byte) (int, error) {
//	return r.ResponseWriter.Write(b)
//}

// WriteHeader writes the response status code and stores it internally
func (r *StatWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func IsStatusError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 400
}
