package userhandler

import (
	_ "embed"
	"fmt"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/schema"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"path"
)

// DbHandler Is an opinionated user management handler for an http server
// it uses uses the same DB user manager to manage user creation, deletion, etc.
// todo add factory to make sure the logger is not nil
// todo rename to something that makes more sense
type DbHandler struct {
	mng *user.DbManager
}

func NewDbHandler(db *gorm.DB, opts user.ManagerOpts) (*DbHandler, error) {

	dbh, err := user.NewDbManager(db, opts)
	if err != nil {
		return nil, err
	}

	hndlr := DbHandler{
		mng: dbh,
	}
	return &hndlr, nil
}

func NewHandler(db *user.DbManager) (*DbHandler, error) {
	hndlr := DbHandler{
		mng: db,
	}
	return &hndlr, nil
}

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var formDecoder = schema.NewDecoder()

type FormPostData struct {
	Name     string
	Pw       string
	Redirect string
}

// tmplData holds the fields passed to the template in the user management forms
type tmplData struct {
	Path     string
	Redirect string
}

//go:embed tmpl/createUserForm.html
var createUserForm string

// CreateUserForm returns a handler function to print a simple create user form
func (h DbHandler) CreateUserForm() func(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.New("name").Parse(createUserForm)
	// todo handle error
	// Error checking elided
	data := tmplData{}

	return func(w http.ResponseWriter, r *http.Request) {
		data.Path = r.RequestURI
		data.Redirect = path.Clean(r.RequestURI + "/..")
		_ = tmpl.Execute(w, data)
	}
}

// CreateUserHandleForm returns a handler function to process user creation post form request
func (h DbHandler) CreateUserHandleForm() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// TODO
			panic(err)
			// Handle error
		}

		var payload FormPostData

		// r.PostForm is a map of our POST form values
		err = formDecoder.Decode(&payload, r.PostForm)
		if err != nil {
			spew.Dump(err)
		}

		err = h.mng.Create(payload.Name, payload.Pw)
		if err != nil {
			spew.Dump(err)
		}

		// todo send path from get request
		http.Redirect(w, r, payload.Redirect, http.StatusSeeOther)

		// todo handle errors
		// todo handle user already exists ?
		// todo add nonce with server storage ?
		// todo handle json payload?

		// redirect depending on the status

	}
}

// UserHandler returns a full-fledged user handling mux
// it is not intended as real use, but rather as an example on how to use all or part of the user manager handler
// it can be hooked in into any other mux like this:
// rootMux.Handle("/top_path/", http.StripPrefix("/top_path", subMux))
func (h DbHandler) UserHandler(base string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(path.Join(base, "/register"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spew.Dump("here")
		if r.Method == http.MethodGet {
			h.CreateUserForm()(w, r)
		} else if r.Method == http.MethodPost {
			h.CreateUserHandleForm()(w, r)
		} else {
			// todo unsuported request
			panic(fmt.Errorf("unsuppreted req"))
		}
	}))

	rootHandler := handlers.SimpleText{
		Text: "user root",
		Links: []handlers.Link{
			{Text: "back", Url: path.Join(base, "../")},
			{Text: "Register", Url: path.Join(base, "/register")},
		},
	}
	mux.Handle("/", &rootHandler)

	return mux
}
