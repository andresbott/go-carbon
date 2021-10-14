package auth

import (
	"net/http"
)

type UserLogin interface {
	AllowLogin(user string, hash string) bool
}

type Basic struct {
	User         UserLogin
	Redirect     string
	RedirectCode int
}

func (auth *Basic) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if ok {
			if auth.User.AllowLogin(username, password) {
				next.ServeHTTP(w, r)
				return
			}
		}

		if auth.Redirect != "" {
			http.Redirect(w, r, auth.Redirect, auth.RedirectCode)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
