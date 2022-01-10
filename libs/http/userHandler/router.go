package userHandler

import (
	"git.andresbott.com/Golang/carbon/libs/http/simpleTextHandler"
	"github.com/gorilla/mux"
	"net/http"
	"path"
)

const (
	CreateUrl = "/create"
	LoginUrl  = "/login"
)

type MuxRouter struct {
	Handler

	SubRoute string
}

// AttachHandlers will take a mux.router and add path handlers for users
func (mr MuxRouter) AttachHandlers(r *mux.Router) {

	// attach all the routes to a subrouter
	if mr.SubRoute != "" {
		r = r.PathPrefix("/user").Subrouter()
	}

	handler := simpleTextHandler.Handler{
		Text: "user root",
		Links: []simpleTextHandler.Link{
			{
				Text: "back",
				Url:  "../",
			},
			{
				Text: "Login",
				Url:  path.Join(mr.SubRoute, LoginUrl),
			},
			{
				Text: "Register",
				Url:  path.Join(mr.SubRoute, CreateUrl),
			},
		},
	}
	r.Path("/").Handler(&handler)
	if mr.SubRoute != "" {
		// in case of a subroute we need to handle empty path for requests to /sub (without tailing slash)
		r.Path("").Handler(&handler)
	}

	// create
	r.Path(CreateUrl).Methods(http.MethodGet).HandlerFunc(mr.CreateUserGetFrom())
	r.Path(CreateUrl).Methods(http.MethodPost).HandlerFunc(mr.CreateUserPostForm())

	// login
	r.Path(LoginUrl).Methods(http.MethodGet).HandlerFunc(mr.LoginUserGetHandler())
	r.Path(LoginUrl).Methods(http.MethodPost).HandlerFunc(mr.LoginUserPostHandler())

}
