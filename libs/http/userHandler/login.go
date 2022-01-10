package userHandler

import (
	_ "embed"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"html/template"
	"net/http"
)

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

		var payload FormPostData

		// r.PostForm is a map of our POST form values
		err = formDecoder.Decode(&payload, r.PostForm)
		if err != nil {
			spew.Dump(err)
		}

		isAllowed := h.Manager.CheckLogin(payload.Name, payload.Pw)
		if isAllowed {
			fmt.Println("Login Allowed")
		} else {
			fmt.Println("Login Denied")
		}

		// todo get redirect path from form post request
		http.Redirect(w, r, payload.Redirect, http.StatusSeeOther)

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
