package router

import (
	"git.andresbott.com/Golang/carbon/app/handlrs"
	"git.andresbott.com/Golang/carbon/internal/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/http/middleware"

	"github.com/gorilla/mux"
	"net/http"
)

//	@title			Carbon Sample API
//	@version		0.1
//	@description	Sample implementation of an API using the carbon framework
// TODO add license to swagger
//	@BasePath		/api/v0

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func apiV0(r *mux.Router, session *auth.SessionMgr, users auth.UserLogin, manager *tasks.Manager) error {

	r.Use(func(handler http.Handler) http.Handler {
		// todo this should reflect prod vs non-prod property
		return middleware.JsonErrMiddleware(handler, false)
	})
	// add users handling to api
	userApi(r, session, users)

	// add tasks api
	tasksApi(r, session, manager)
	return nil
}

func tasksApi(r *mux.Router, session *auth.SessionMgr, manager *tasks.Manager) {
	pageHandler := handlers.SimpleText{
		Text: "Page protected by session auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}
	r.Use(session.Middleware)

	th := handlrs.TaskHandler{
		TaskManager: manager,
	}

	ProtectedPage := session.Middleware(&pageHandler)
	// GET
	r.Path("/tasks").Handler(ProtectedPage)
	// PUT
	r.Path("/task").Methods(http.MethodPut).Handler(th.Create())
	// GET | PUT | DELETE
	r.Path("/task/{ID}").Handler(ProtectedPage)
}

func userApi(apiRoute *mux.Router, session *auth.SessionMgr, users auth.UserLogin) {
	userLogin(apiRoute, session, users)
	userLogout(apiRoute, session)
	userStatus(apiRoute, session)
	userOptions(apiRoute, session)
}

// userLogin
//
//	@Summary		Login a user
//	@Description	Handles the user login process
//	@Tags			User
//	@Produce		json
//	@Param			UserData	body		handlrs.loginData	true	"user login payload"
//	@Success		200			{object}	handlrs.userStatus
//	@Router			/user/login [post]
func userLogin(apiRoute *mux.Router, session *auth.SessionMgr, users auth.UserLogin) {
	apiRoute.Path("/user/login").Methods(http.MethodPost).Handler(handlrs.UserLoginHandler(session, users))
	apiRoute.Path("/user/login").Methods(http.MethodOptions).Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))
	apiRoute.Path("/user/login").Handler(handlers.StatusErr(http.StatusMethodNotAllowed))
}

// userLogout
//
//	@Summary		Logout a user
//	@Description	Logs out the current user based on the session cookie
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	handlrs.userStatus
//	@Router			/user/logout [get]
//	@Router			/user/logout [put]
//	@Router			/user/logout [post]
func userLogout(apiRoute *mux.Router, session *auth.SessionMgr) {
	apiRoute.Path("/user/logout").Handler(handlrs.UserLogoutHandler(session))
}

// userOptions
//
//	@Summary		Get user options
//	@Tags			User
//	@Description	Get options specific to the currently logged-in user based on the session cookie
//	@Produce		json
//	@Success		501	{object}	middleware.jsonErr
//	@Failure		405	{object}	middleware.jsonErr
//	@Failure		500	{object}	middleware.jsonErr
//	@Router			/user/options [get]
func userOptions(apiRoute *mux.Router, session *auth.SessionMgr) {
	apiRoute.Path("/user/options").Methods(http.MethodGet).Handler(handlers.StatusErr(http.StatusNotImplemented))
	apiRoute.Path("/user/options").Handler(handlers.StatusErr(http.StatusMethodNotAllowed))
}

// userStatus
//
//	@Tags			User
//	@Summary		Get user status
//	@Description	Show the satus information about the current user
//	@Produce		json
//	@Success		200	{object}	handlrs.userStatus
//	@Router			/user/status [get]
func userStatus(apiRoute *mux.Router, session *auth.SessionMgr) {
	apiRoute.Path("/user/status").Methods(http.MethodGet).Handler(handlrs.UserStatusHandler(session))
	apiRoute.Path("/user/status").Handler(handlers.StatusErr(http.StatusMethodNotAllowed))
}
