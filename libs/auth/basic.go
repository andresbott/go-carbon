package auth

import (
	"git.andresbott.com/Golang/carbon/libs/log"
	"net/http"
)

type UserLogin interface {
	AllowLogin(user string, password string) bool
}

type Basic struct {
	User         UserLogin
	Redirect     string
	RedirectCode int
	Logger       log.LeveledLogger
}

func (auth *Basic) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth.Logger.Info("auth")
		username, password, ok := r.BasicAuth()

		if ok {
			if auth.User.AllowLogin(username, password) {
				auth.Logger.Info("ok")
				next.ServeHTTP(w, r)
				return
			}
		}

		if auth.Redirect != "" {
			http.Redirect(w, r, auth.Redirect, auth.RedirectCode)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Access to the staging site", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
