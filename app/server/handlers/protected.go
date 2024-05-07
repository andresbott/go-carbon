package handlers

import (
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"net/http"
)

func fixedBasicAuthHandler() http.Handler {

	page := handlers.SimpleText{
		Text: "Page protected by basic auth",
		Links: []handlers.Link{
			{
				Text: "back to root",
				Url:  "../",
			},
		},
	}
	return &page
}

func dbBasicAuthHandler() http.Handler {

	page := handlers.SimpleText{
		Text: "Page protected by basic auth with users in a DB",
		Links: []handlers.Link{
			{
				Text: "back to root",
				Url:  "../",
			},
		},
	}

	return &page
}
func cookieProtectedContent() http.Handler {

	page := handlers.SimpleText{
		Text: "Page protected by cookie auth",
		Links: []handlers.Link{
			{
				Text: "back to root",
				Url:  "../",
			},
		},
	}
	return &page
}
