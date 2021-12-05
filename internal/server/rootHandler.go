package server

import (
	"git.andresbott.com/Golang/carbon/libs/http/handlers/simpleText"
	"git.andresbott.com/Golang/carbon/libs/http/handlers/userhandler"
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
func newRootHandler(l log.LeveledLogger, db *gorm.DB) *rootHandler {

	r := mux.NewRouter()

	// add logging middleware
	r.Use(func(handler http.Handler) http.Handler {
		return log.LoggingMiddleware(handler, l)
	})

	// root page
	// --------------------------
	rootPage := simpleText.Handler{
		Text: "root page",
		Links: []simpleText.Link{
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

	r.Path("/").Handler(&rootPage)

	// user handling
	// --------------------------
	userManager, err := user.NewManager(db, user.ManagerOpts{
		BcryptDifficulty: 4,
	})
	if err != nil {
		panic(err) // it is ok to panic during startup
	}

	userHandler := userhandler.Handler{
		Manager:  userManager,
		SubRoute: "/user",
	}
	userHandler.AttachHandlers(r)

	// page protected by basic auth
	// --------------------------
	basicAuthPage := simpleText.Handler{
		Text: "Page protected by basic auth",
		Links: []simpleText.Link{
			{
				Text: "back to root",
				Url:  "../",
			},
		},
	}

	//basicAuth := auth.Basic{
	//	User:         dummyUser{},
	//	Redirect:     "",
	//	RedirectCode: 302,
	//	Logger:       l,
	//}
	r.Path("/basic").Handler(userHandler.GetBasicLoginMiddleware(&basicAuthPage))

	handlr := rootHandler{
		router: r,
	}

	return &handlr
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
