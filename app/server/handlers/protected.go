package handlers

import (
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"github.com/gorilla/mux"
)

func basiAuthProtectPage(r *mux.Router) {

	page := handlers.SimpleText{
		Text: "Page protected by basic auth",
		Links: []handlers.Link{
			{
				Text: "back to root",
				Url:  "../",
			},
		},
	}

	users := auth.FixedUsers{
		Users: map[string]string{
			"demo": "demo",
		},
	}
	authProtected := auth.Basic{
		User: users,
	}

	r.Path("/basic").Handler(authProtected.Middleware(&page))
}
