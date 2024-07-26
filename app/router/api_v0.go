package router

import (
	"git.andresbott.com/Golang/carbon/app/handlrs"
	"git.andresbott.com/Golang/carbon/internal/model/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/http/middleware"

	"github.com/gorilla/mux"
	"net/http"
)

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

// tasksApi
func tasksApi(r *mux.Router, session *auth.SessionMgr, manager *tasks.Manager) {
	r.Use(session.Middleware)
	th := handlrs.TaskHandler{
		TaskManager: manager,
	}
	r.Path("/tasks").Methods(http.MethodGet).Handler(th.List())

	r.Path("/task").Methods(http.MethodPost).Handler(th.Create())
	r.Path("/task/{ID}").Methods(http.MethodGet).Handler(th.Read())
	r.Path("/task/{ID}").Methods(http.MethodDelete).Handler(th.Delete())
	r.Path("/task/{ID}").Methods(http.MethodPut).Handler(th.Update())
}

func userApi(apiRoute *mux.Router, session *auth.SessionMgr, users auth.UserLogin) {

	//  LOGIN
	apiRoute.Path("/user/login").Methods(http.MethodPost).Handler(handlrs.UserLoginHandler(session, users))
	apiRoute.Path("/user/login").Methods(http.MethodOptions).Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))
	apiRoute.Path("/user/login").Handler(handlers.StatusErr(http.StatusMethodNotAllowed))

	// LOGOUT
	apiRoute.Path("/user/logout").Handler(handlrs.UserLogoutHandler(session))

	// STATUS
	apiRoute.Path("/user/status").Methods(http.MethodGet).Handler(handlrs.UserStatusHandler(session))
	apiRoute.Path("/user/status").Handler(handlers.StatusErr(http.StatusMethodNotAllowed))

	// OPTIONS
	apiRoute.Path("/user/options").Methods(http.MethodGet).Handler(handlers.StatusErr(http.StatusNotImplemented))
	apiRoute.Path("/user/options").Handler(handlers.StatusErr(http.StatusMethodNotAllowed))
}
