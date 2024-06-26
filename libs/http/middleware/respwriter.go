package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

// StatWriter is a wrapper to a httpResponse writer that allows to intercept and
// extract the status code that the upstream code has defined
type StatWriter struct {
	http.ResponseWriter
	statusCode    int
	InterceptBody bool // write a limited amount of chars into a buffer in case a non 200 code
	buf           limBuf
}

// NewWriter will return a pointer to a response writer
func NewWriter(w http.ResponseWriter) *StatWriter {
	// WriteHeader(int) is not called if the response is 200 (implicit response code) so it needs to be the default
	return &StatWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		buf: limBuf{
			cap: 1000,
			buf: &bytes.Buffer{},
		},
	}
}

func (r *StatWriter) StatusCode() int {
	return r.statusCode
}

func (r *StatWriter) StatusCodeStr() string {
	return strconv.Itoa(r.statusCode)
}

// Write returns underlying Write result, while counting data size
func (r *StatWriter) Write(b []byte) (int, error) {
	if r.InterceptBody && r.statusCode != 200 {
		return r.buf.Write(b)
	}
	return r.ResponseWriter.Write(b)
}

// WriteHeader writes the response status code and stores it internally
func (r *StatWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func IsStatusError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 400
}

// limited buffer than only accepts up to certain size
type limBuf struct {
	cap int
	buf *bytes.Buffer
}

func (b *limBuf) Write(p []byte) (n int, err error) {
	if len(p)+b.buf.Len() >= b.cap {
		return len(p), fmt.Errorf("buf limit reached")
	} else {
		b.buf.Write(p)
	}
	return len(p), nil
}
