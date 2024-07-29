package router

import (
	"github.com/andresbott/go-carbon/libs/auth"
	"github.com/andresbott/go-carbon/libs/http/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func demo(r *mux.Router) error {
	demoPage := handlers.SimpleText{
		Text: "Demo root page",
		Links: []handlers.Link{
			{Text: "Basic auth protected (demo:demo)", Url: "/basic"},
			{Text: "SPA: load the SPA", Url: "/spa"},
			{Text: "session Authentication protected page", Url: "/session"},
			//{
			//	Text: "session based login (demo:demo)",
			//	Url:  handlers2.sessionLogin,
			//},
			{Text: "User handling", Url: "/user"},
			{Text: "Json API", Child: []handlers.Link{
				{Text: "User Status", Url: "/api/v0/user/status"},
				{Text: "User options", Url: "/api/v0/user/options"},
			}},
			{Text: "Observability", Child: []handlers.Link{
				{Text: "metrics", Url: "http://localhost:9090/metrics"},
			}},
		},
	}
	r.Path("/").Handler(demoPage)

	// user management
	// --------------------------
	// db managed users
	//sampleDbUser, err := sampleUserManager()
	//if err != nil {
	//	return nil, err
	//}
	//
	//userDbHandler, err := userhandler.NewHandler(sampleDbUser)
	//if err != nil {
	//	return nil, err
	//}
	//r.PathPrefix("/user").Handler(userDbHandler.UserHandler("/user"))

	return nil

}

func basicAuthProtected(r *mux.Router, users auth.UserLogin) error {

	demoPage := handlers.SimpleText{
		Text: "Basic auth protected page",
		Links: []handlers.Link{
			{
				Text: "Back",
				Url:  "/",
			},
		},
	}
	bAuth := auth.Basic{
		User: users,
	}

	// use the middleware to protect the page
	r.Use(func(handler http.Handler) http.Handler {
		return bAuth.Middleware(handler)
	})

	r.Path("").Handler(demoPage)

	return nil
}

// const SessionLogin = "/session-login"

func SessionProtected(r *mux.Router, session *auth.SessionMgr) error {
	pageHandler := handlers.SimpleText{
		Text: "Page protected by session auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}

	ProtectedPage := session.Middleware(&pageHandler)
	r.Path("/session").Handler(ProtectedPage)

	return nil
}
