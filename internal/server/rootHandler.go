package server

import (
	"git.andresbott.com/Golang/carbon/libs/http/simpleTextHandler"
	"git.andresbott.com/Golang/carbon/libs/http/userHandler"
	"git.andresbott.com/Golang/carbon/libs/log"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type rootHandler struct {
	router *mux.Router
}

// newRootHandler generates the main url router handler to be used in the server
func newRootHandler(l log.LeveledStructuredLogger, db *gorm.DB) *rootHandler {

	r := mux.NewRouter()

	// add logging middleware
	r.Use(func(handler http.Handler) http.Handler {
		return log.LoggingMiddleware(handler, l)
	})

	// root page
	// --------------------------
	rootPage := simpleTextHandler.Handler{
		Text: "root page",
		Links: []simpleTextHandler.Link{
			{
				Text: "Basic auth protected",
				Url:  "/basic",
			},
			{
				Text: "User handling",
				Url:  "/user",
			},
		},
	}

	// user handling
	// --------------------------
	userManager, err := user.NewManager(db, user.ManagerOpts{
		BcryptDifficulty: 4,
	})
	if err != nil {
		panic(err) // it is ok to panic during startup
	}

	usrHndlr := userHandler.MuxRouter{
		Handler:  userHandler.Handler{Manager: userManager},
		SubRoute: "/user",
	}
	usrHndlr.AttachHandlers(r)

	basiAuthProtectPage(r, usrHndlr.Handler)
	r.Path("/").Handler(&rootPage)

	return &rootHandler{
		router: r,
	}
}

// add a basic auth protected page
func basiAuthProtectPage(r *mux.Router, usrHndlr userHandler.Handler) {

	basicAuthPage := simpleTextHandler.Handler{
		Text: "Page protected by basic auth",
		Links: []simpleTextHandler.Link{
			{
				Text: "back to root",
				Url:  "../",
			},
		},
	}
	r.Path("/basic").Handler(usrHndlr.GetBasicLoginMiddleware(&basicAuthPage))
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
