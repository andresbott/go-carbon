package routes

import (
	"git.andresbott.com/Golang/carbon/app/server/handlers"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/middleware"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
	"time"
)

func ApiV0(r *mux.Router) error {

	apiRoute := r.PathPrefix("/api/v0").Subrouter()
	r.Use(func(handler http.Handler) http.Handler {
		return middleware.JsonErrMiddleware(handler)
	})

	store, err := auth.FsStore("", securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	if err != nil {
		return err
	}
	// create an instance of session auth
	sessionAuth, _ := auth.NewSessionMgr(auth.SessionCfg{
		Store:         store,
		SessionDur:    time.Hour,       // time the user is logged in
		MaxSessionDur: 24 * time.Hour,  // time after the user is forced to re-login anyway
		MinWriteSpace: 2 * time.Minute, // throttle write operations on the session
	})

	users := user.StaticUsers{
		Users: map[string]string{
			"demo": "demo",
		},
	}
	apiRoute.Path("/user/login").Methods(http.MethodPost).Handler(handlers.UserLoginHandler(sessionAuth, users))
	apiRoute.Path("/user/login").Handler(middleware.ErrorHandler(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed))

	return nil
}
