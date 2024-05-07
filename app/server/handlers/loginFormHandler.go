package handlers

import (
	_ "embed"
	"github.com/gorilla/schema"
	"html/template"
	"net/http"
	"path"
)

//go:embed tmpl/loginForm.html
var loginForm string

// tmplData holds the fields passed to the template in the user management forms
type tmplData struct {
	Path     string
	Redirect string
}

type LoginFormData struct {
	Name     string
	Pw       string
	Redirect string
}

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

func LoginForm() func(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.New("name").Parse(loginForm)
	// todo handle error
	// Error checking elided
	data := tmplData{}

	return func(w http.ResponseWriter, r *http.Request) {
		data.Path = r.RequestURI
		data.Redirect = path.Clean(r.RequestURI + "/..")
		_ = tmpl.Execute(w, data)
	}
}
