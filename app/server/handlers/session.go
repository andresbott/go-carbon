package handlers

import (
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/gorilla/mux"
	"net/http"
)

const sessionLogin = "/session-login"
const sessionContent = "/session"

func sessionAuthentication(r *mux.Router, demoUsers user.StaticUsers) error {
	// session based auth
	// --------------------------
	hashKey := []byte("oach9iu2uavahcheephi4FahzaeNge8yeecie4jee9rah9ahrah6tithai7Oow5U")
	blockKey := []byte("eeth3oon5eewifaogeibieShey5eiJ0E")

	cookieStore, err := auth.CookieStore(hashKey, blockKey)
	if err != nil {
		return err
	}

	sessionAuth, err := auth.NewSessionMgr(auth.SessionCfg{
		User:  demoUsers,
		Store: cookieStore,
	})
	if err != nil {
		return err
	}

	cookieProtectedPageHandler := handlers.SimpleText{
		Text: "Page protected by session auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}
	cookieProtected := sessionAuth.Middleware(&cookieProtectedPageHandler)
	r.Path(sessionContent).Handler(cookieProtected)
	// handle the post request
	r.PathPrefix(sessionLogin).Methods(http.MethodPost).Handler(sessionAuth.FormAuthHandler())

	// render the form
	loginFormHandlr := handlers.TemplateWithRequest(loginForm, func(r *http.Request) map[string]interface{} {
		payload := map[string]interface{}{}
		payload["Path"] = sessionLogin
		payload["Redirect"] = sessionContent
		return payload
	})
	r.PathPrefix(sessionLogin).HandlerFunc(loginFormHandlr)
	return nil
}
