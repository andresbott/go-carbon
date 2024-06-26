package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonErr struct {
	Error string
	Code  int
}

func JsonErrorHandler(err string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, err, code)
	})
}
func JsonErrMethodNotAllowed() http.Handler {
	return JsonErrorHandler("method not allowed", http.StatusMethodNotAllowed)
}

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

func ErrorHandler(err string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err, code)
	})
}
