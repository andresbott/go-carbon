package v0

import (
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/gorilla/mux"
)

func basicAuth(r *mux.Router, demoUsers user.StaticUsers) {
	// Basic auth protected path
	// --------------------------

	fixedAuth := auth.Basic{
		User: demoUsers,
	}
	fixedAuthPageHandlr := handlers.SimpleText{
		Text: "Page protected by basic auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}
	fixedProtectedPath := fixedAuth.Middleware(&fixedAuthPageHandlr)
	r.Path("/basic").Handler(fixedProtectedPath)
}
