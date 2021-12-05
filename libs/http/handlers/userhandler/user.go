package userhandler

import (
	_ "embed"
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/http/handlers/simpleText"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"html/template"
	"net/http"
	"path"
)

const (
	CreateUrl = "/create"
	LoginUrl  = "/login"
)

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

// Handler Is an opinionated user management endpoint for an http server
// it uses User manager interface to manage user creation, deletion, etc.
// todo add factory to make sure the logger is not nil
// todo rename to something that makes more sense
type Handler struct {
	Manager UserManager

	SubRoute string
}

type UserManager interface {
	Create(id string, pw string) error
	CheckLogin(id string, pw string) bool
}

// todo login type: basic, form-post, json-post, token, oauth

// AttachHandlers will take a mux.router and add path handlers for users
func (h Handler) AttachHandlers(r *mux.Router) {

	// attach all the routes to a subrouter
	if h.SubRoute != "" {
		r = r.PathPrefix("/user").Subrouter()
	}

	handler := simpleText.Handler{
		Text: "user root",
		Links: []simpleText.Link{
			{
				Text: "back",
				Url:  "../",
			},
			{
				Text: "Login",
				Url:  path.Join(h.SubRoute, LoginUrl),
			},
			{
				Text: "Register",
				Url:  path.Join(h.SubRoute, CreateUrl),
			},
		},
	}
	r.Path("/").Handler(&handler)
	if h.SubRoute != "" {
		// in case of a subroute we need to handle empty path for requests to /sub (without tailing slash)
		r.Path("").Handler(&handler)
	}

	// create
	r.Path(CreateUrl).Methods(http.MethodGet).HandlerFunc(h.CreateUserGetHandler())
	r.Path(CreateUrl).Methods(http.MethodPost).HandlerFunc(h.CreateUserPostHandler())

	// login
	r.Path(LoginUrl).Methods(http.MethodGet).HandlerFunc(h.LoginUserGetHandler())
	r.Path(LoginUrl).Methods(http.MethodPost).HandlerFunc(h.LoginUserPostHandler())

}

type tmplData struct {
	Path string
}

//go:embed tmpl/createUserForm.html
var createUserForm string

//CreateUserGetHandler returns a handler function to process user creation post requests
func (h Handler) CreateUserGetHandler() func(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.New("name").Parse(createUserForm)
	// Error checking elided
	data := tmplData{
		Path: "",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data.Path = r.RequestURI
		_ = tmpl.Execute(w, data)
		return
	}
}

type User struct {
	Name string
	Pw   string
}

//CreateUserPostHandler returns a handler function to process user creation post requests
func (h Handler) CreateUserPostHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// Handle error
		}

		var usr User

		// r.PostForm is a map of our POST form values
		err = decoder.Decode(&usr, r.PostForm)
		if err != nil {
			spew.Dump(err)
		}

		err = h.Manager.Create(usr.Name, usr.Pw)
		if err != nil {
			spew.Dump(err)
		}
		fmt.Println("POST request")
		http.Redirect(w, r, h.SubRoute, http.StatusSeeOther)

		return
		// todo handle errors
		// todo handle user already exists ?
		// todo add nonce with server storage ?
		// todo handle json payload?

		// redirect depending on the status

	}
}

//go:embed tmpl/loginForm.html
var loginForm string

//LoginUserGetHandler returns a handler function to process user login
func (h Handler) LoginUserGetHandler() func(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.New("name").Parse(loginForm)
	// Error checking elided
	data := tmplData{
		Path: "",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data.Path = r.RequestURI
		_ = tmpl.Execute(w, data)
		return
	}
}

//LoginUserPostHandler returns a handler function to process user creation post requests
func (h Handler) LoginUserPostHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// Handle error
		}

		var usr User

		// r.PostForm is a map of our POST form values
		err = decoder.Decode(&usr, r.PostForm)
		if err != nil {
			spew.Dump(err)
		}

		isAllowed := h.Manager.CheckLogin(usr.Name, usr.Pw)
		if isAllowed {
			fmt.Println("Login Allowed")
		} else {
			fmt.Println("Login Denied")
		}

		http.Redirect(w, r, h.SubRoute, http.StatusSeeOther)

		return
		// todo handle errors
		// todo handle user already exists ?
		// todo add nonce with server storage ?
		// todo handle json payload?

		// redirect depending on the status

	}
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
