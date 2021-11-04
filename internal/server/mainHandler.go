package server

import (
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/log"
	"github.com/gorilla/mux"
	"net/http"
)

type mainHandler struct {
	router *mux.Router
	Logger log.LeveledLogger
}

type dummyUser struct {
}

func (st dummyUser) AllowLogin(user string, hash string) bool {
	if user == "admin" && hash == "admin" {
		return true
	}
	return false
}

func newMainHandler(l log.LeveledLogger) *mainHandler {

	r := mux.NewRouter()
	// root page
	r.Path("/").Handler(&textHandler{
		Text: "root page",
		Links: map[string]string{
			"basic": "/basic",
		},
		Logger: l,
	})

	// page protected by basic auth

	basicAuth := auth.Basic{
		User:         dummyUser{},
		Redirect:     "",
		RedirectCode: 302,
		Logger:       l,
	}

	basiHandler := textHandler{
		Text: "basic auth",
		Links: map[string]string{
			"root": "../",
		},
		Logger: l,
	}

	r.Path("/basic").Handler(basicAuth.Middleware(&basiHandler))

	handler := mainHandler{
		router: r,
		Logger: l,
	}
	return &handler
}

func (h *mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//s.logger.Debug(fmt.Sprintf("serving request on url: %s method: %v\n", r.URL, r.Method))
	h.router.ServeHTTP(w, r)
}
