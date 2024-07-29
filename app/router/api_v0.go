package router

import (
	"github.com/andresbott/go-carbon/app/handlrs"
	"github.com/andresbott/go-carbon/internal/model/tasks"
	"github.com/andresbott/go-carbon/libs/auth"
	"github.com/andresbott/go-carbon/libs/http/handlers"
	"github.com/andresbott/go-carbon/libs/http/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func apiV0(r *mux.Router, session *auth.SessionMgr, users auth.UserLogin, manager *tasks.Manager) error {
	// todo this should reflect prod vs non-prod property
	genericErrorMessage := false

	// this sub router does NOT enforce authentication
	openSubRoute := r.PathPrefix("/api/v0").Subrouter()
	openSubRoute.Use(func(handler http.Handler) http.Handler {
		return middleware.JsonErrMiddleware(handler, genericErrorMessage)
	})
	// add users handling to api
	userApi(openSubRoute, session, users)

	// this sub router does enforce authentication
	protected := r.PathPrefix("/api/v0").Subrouter()
	protected.Use(func(handler http.Handler) http.Handler {
		return middleware.JsonErrMiddleware(handler, genericErrorMessage)
	}, session.Middleware)

	// add tasks api
	tasksApi(protected, manager)
	return nil
}

// tasksApi
func tasksApi(r *mux.Router, manager *tasks.Manager) {
	th := handlrs.TaskHandler{TaskManager: manager}
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
