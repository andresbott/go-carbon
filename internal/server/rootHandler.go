package server

import (
	"git.andresbott.com/Golang/carbon/internal/server/textHandler"
	userHandler "git.andresbott.com/Golang/carbon/internal/server/user"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/log"
	"git.andresbott.com/Golang/carbon/libs/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type dummyUser struct {
}

func (st dummyUser) AllowLogin(user string, hash string) bool {
	if user == "admin" && hash == "admin" {
		return true
	}
	return false
}

type rootHandler struct {
	router *mux.Router
}

// redirectHandler redirects the request to the desired location
func redirectHandler(url string, permanent bool) func(w http.ResponseWriter, r *http.Request) {
	code := http.StatusTemporaryRedirect
	if permanent {
		code = http.StatusPermanentRedirect
	}
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, code)
	}
}

// newRootHandler generates the main url router handler to be used in the server
func newRootHandler(l log.LeveledLogger) *rootHandler {

	r := mux.NewRouter()

	// add logging middleware
	r.Use(func(handler http.Handler) http.Handler {
		return middleware.LoggingMiddleware(handler, l)
	})

	// root page
	// --------------------------
	rootPage := textHandler.Handler{
		Text: "root page",
		Links: map[string]string{
			"basic": "/basic",
			"user":  "/user",
		},
	}

	r.Path("/").Handler(middleware.LoggingMiddleware(&rootPage, l))

	// user handling
	// --------------------------
	r.Path("/user").HandlerFunc(redirectHandler("/user/", true))
	userHandler.UserRoutes(r.PathPrefix("/user").Subrouter())

	// page protected by basic auth
	// --------------------------
	basicAuth := auth.Basic{
		User:         dummyUser{},
		Redirect:     "",
		RedirectCode: 302,
		Logger:       l,
	}

	basicAuthPage := textHandler.Handler{
		Text: "basic auth",
		Links: map[string]string{
			"root": "../",
		},
	}
	r.Path("/basic").Handler(basicAuth.Middleware(&basicAuthPage))

	handlr := rootHandler{
		router: r,
	}

	return &handlr
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//s.logger.Debug(fmt.Sprintf("serving request on url: %s method: %v\n", r.URL, r.Method))
	h.router.ServeHTTP(w, r)
}
