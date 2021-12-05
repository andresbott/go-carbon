package userhandler

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

func (h Handler) GetBasicLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()

		spew.Dump(username, password)
		// basic auth is present
		if ok {
			if h.Manager.CheckLogin(username, password) {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
