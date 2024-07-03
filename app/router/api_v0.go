package router

import (
	"git.andresbott.com/Golang/carbon/app/handlrs"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func apiV0(r *mux.Router, session *auth.SessionMgr, users auth.UserLogin) error {

	r.Use(func(handler http.Handler) http.Handler {
		// todo this should reflect prod vs non-prod property
		return middleware.JsonErrMiddleware(handler, false)
	})
	// add users handling to api
	err := apiV0User(r, session, users)
	if err != nil {
		return err
	}
	return nil
}

func apiV0User(apiRoute *mux.Router, session *auth.SessionMgr, users auth.UserLogin) error {
	apiRoute.Path("/user/status").Methods(http.MethodGet).Handler(handlrs.UserStatusHandler(session))
	apiRoute.Path("/user/status").Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}))

	apiRoute.Path("/user/options").Methods(http.MethodGet).Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}))
	apiRoute.Path("/user/options").Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}))

	apiRoute.Path("/user/login").Methods(http.MethodPost).Handler(handlrs.UserLoginHandler(session, users))
	apiRoute.Path("/user/login").Methods(http.MethodOptions).Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		// TODO either move to middleware or other type of control mechanisms
		headers := writer.Header()
		headers.Add("Access-Control-Allow-Origin", "*")
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
		headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
	}))
	apiRoute.Path("/user/login").Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}))

	return nil
}
