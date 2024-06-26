package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonErr struct {
	Error string
	Code  int
}

//	func JsonErrorHandler(err string, code int) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			jsonError(w, err, code)
//		})
//	}
//
//	func JsonErrMethodNotAllowed() http.Handler {
//		return JsonErrorHandler("method not allowed", http.StatusMethodNotAllowed)
//	}
func jsonError(w http.ResponseWriter, error string, code int) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	payload := jsonErr{
		Error: error,
		Code:  code,
	}
	byteErr, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, string(byteErr))
}

//
//func ErrorHandler(err string, code int) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		http.Error(w, err, code)
//	})
//}

HERE + test
// TODO configuration to return verbose error or generic
// todo drop the original response writer in case of error
// TODO unify in the request logger midlleware
// are there any response codes other than 200 that contain usable body?
func JsonErrMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respWriter := NewWriter(w)
		respWriter.InterceptBody = true
		next.ServeHTTP(respWriter, r)

		code := respWriter.StatusCode()
		if IsStatusError(code) {
			jsonError(w, http.StatusText(code), code)
		}
	})
}
