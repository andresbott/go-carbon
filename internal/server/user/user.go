package userHandler

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/internal/server/textHandler"
	"git.andresbott.com/Golang/carbon/libs/log"
	"github.com/gorilla/mux"
	"net/http"
)

// UserRoutes populates multiple handlers in a single sub route
func UserRoutes(r *mux.Router) {
	// create
	r.Path("/").Handler(&textHandler.Handler{
		Text:  "user root",
		Links: nil,
	})
}

type UserHandler struct {
	logger log.LeveledStructuredLogger
}

type CreateUserHandler struct {
	UserHandler
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", " text/html")
		fmt.Fprintf(w, "Here goes the form to register")
		return
	}

	if r.Method == http.MethodPost {
		h.logger.InfoW("Handle create user post operation")
		// redirect depending on the status
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

type Handler struct {
}

//
// Create -> POST to user/
// email verification >
// Read -> Get to User/profile
// 			if not logged in-> redirect to Login
// Update -> POST to user/profile (using userId in form)
// Delete ->
// login -> POST to user/login
// logout -> Post to user/logout (if possible)

// admin/user/Create?
// admin/user/manage ( list, bulk deactivate)
