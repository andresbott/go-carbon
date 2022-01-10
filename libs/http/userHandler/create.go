package userHandler

import (
	_ "embed"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/schema"
	"html/template"
	"net/http"
	"path"
)

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

// Handler Is an opinionated user management endpoint for an http server
// it uses User manager interface to manage user creation, deletion, etc.
// todo add factory to make sure the logger is not nil
// todo rename to something that makes more sense
type Handler struct {
	Manager UserManager
}

type UserManager interface {
	Create(id string, pw string) error
	CheckLogin(id string, pw string) bool
}

// todo login type: basic, form-post, json-post, token, oauth

// tmplData holds the fields passed to the template in the user management forms
type tmplData struct {
	Path     string
	Redirect string
}

//go:embed tmpl/createUserForm.html
var createUserForm string

//CreateUserGetFrom returns a handler function to process user creation get requests
func (h Handler) CreateUserGetFrom() func(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.New("name").Parse(createUserForm)
	// Error checking elided
	data := tmplData{}

	// todo use nonce

	return func(w http.ResponseWriter, r *http.Request) {
		data.Path = r.RequestURI
		data.Redirect = path.Clean(r.RequestURI + "/..")
		_ = tmpl.Execute(w, data)
		return
	}
}

type FormPostData struct {
	Name     string
	Pw       string
	Redirect string
}

//CreateUserPostForm returns a handler function to process user creation post form request
func (h Handler) CreateUserPostForm() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// Handle error
		}

		var payload FormPostData

		// r.PostForm is a map of our POST form values
		err = formDecoder.Decode(&payload, r.PostForm)
		if err != nil {
			spew.Dump(err)
		}

		err = h.Manager.Create(payload.Name, payload.Pw)
		if err != nil {
			spew.Dump(err)
		}

		// todo send path from get request
		http.Redirect(w, r, payload.Redirect, http.StatusSeeOther)

		return
		// todo handle errors
		// todo handle user already exists ?
		// todo add nonce with server storage ?
		// todo handle json payload?

		// redirect depending on the status

	}
}

func (h Handler) CreateUserPostJson() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo handle post json request
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
