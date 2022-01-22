package responseStatusCode

import (
	"net/http"
	"strconv"
)

// ResWriter is a wrapper to a httpResponse writer that allows to intercept and
// extract the status code that the upstream code has defined
type ResWriter struct {
	http.ResponseWriter
	statusCode int
}

// New will return a pointer to a response writer
func New(w http.ResponseWriter) *ResWriter {
	// WriteHeader(int) is not called if the response is 200 (implicit response code) so it needs to be the default
	return &ResWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *ResWriter) StatusCode() int {
	return r.statusCode
}

func (r *ResWriter) StatusCodeStr() string {
	return strconv.Itoa(r.statusCode)
}

//// Write returns underlying Write result, while counting data size
//func (r *ResWriter) Write(b []byte) (int, error) {
//	return r.ResponseWriter.Write(b)
//}

// WriteHeader writes the response status code and stores it internally
func (r *ResWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func IsStatusError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 400
}
