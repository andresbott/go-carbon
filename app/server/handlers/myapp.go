package handlers

import (
	_ "embed"
	"git.andresbott.com/Golang/carbon/internal/http/userhandler"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/log/zero"
	"git.andresbott.com/Golang/carbon/libs/prometheus"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

type MyAppHandler struct {
	router *mux.Router
}

func (h *MyAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

//go:embed tmpl/loginForm.html
var loginForm string

// TODO split the App handler in busines functions

// NewAppHandler generates the main url router handler to be used in the server
func NewAppHandler(l *zerolog.Logger, db *gorm.DB) (*MyAppHandler, error) {

	r := mux.NewRouter()

	// add logging middleware
	r.Use(func(handler http.Handler) http.Handler {
		return zero.LoggingMiddleware(handler, l)
	})

	promMiddle := prometheus.NewMiddleware(prometheus.Cfg{
		MetricPrefix: "myApp",
	})
	r.Use(func(handler http.Handler) http.Handler {
		return promMiddle.Handler(handler)
	})

	demoUsers := user.StaticUsers{
		Users: map[string]string{
			"demo": "demo",
		},
	}

	// Basic auth protected path
	// --------------------------
	fixedAuth := auth.Basic{
		User: demoUsers,
	}
	fixedAuthPageHandlr := handlers.SimpleText{
		Text: "Page protected by basic auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}
	fixedProtectedPath := fixedAuth.Middleware(&fixedAuthPageHandlr)

	r.Path("/basic").Handler(fixedProtectedPath)

	// Basic auth protected path but with demoUsers managed by an in-memory DB
	// --------------------------

	sampleDbUser, err := sampleUserManager()
	if err != nil {
		return nil, err
	}
	authProtected := auth.Basic{
		User: sampleDbUser,
	}
	basicAuthPageHandlr := handlers.SimpleText{
		Text: "Page protected by basic auth with users in a DB",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}
	dbProtectedPath := authProtected.Middleware(&basicAuthPageHandlr)

	r.Path("/basic-auth-db").Handler(dbProtectedPath)

	// Cookie based login and protected content
	// --------------------------
	cookieAuth, err := auth.NewCookieAuth(
		auth.CookieCfg{
			User:         demoUsers,
			Redirect:     "/cookie-login",
			RedirectCode: http.StatusTemporaryRedirect,
			CookieName:   "",
			HashKey:      []byte("banana"),
			BlockKey:     []byte("thahsh0fee4Zae3taizieN9goquie4ze"),
		},
	)
	if err != nil {
		return nil, err
	}

	cookieProtectedPageHandler := handlers.SimpleText{
		Text: "Page protected by cookie auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}
	cookieProtected := cookieAuth.Middleware(&cookieProtectedPageHandler)
	r.Path("/cookie").Handler(cookieProtected)
	// handle the post request
	r.PathPrefix("/cookie-login").Methods(http.MethodPost).Handler(cookieAuth.FormAuthHandler())

	// render the form
	loginFormHandlr := handlers.TmplWithReq(loginForm, func(r *http.Request) map[string]interface{} {
		payload := map[string]interface{}{}
		payload["Path"] = r.RequestURI
		payload["Redirect"] = path.Clean(r.RequestURI + "/..")
		return payload
	})
	r.PathPrefix("/cookie-login").HandlerFunc(loginFormHandlr)

	// user management
	// --------------------------
	userDbHandler, err := userhandler.NewHandler(sampleDbUser)
	if err != nil {
		return nil, err
	}
	r.PathPrefix("/user").Handler(userDbHandler.UserHandler("/user"))

	// root page
	// --------------------------
	rootPage := handlers.SimpleText{
		Text: "root page",
		Links: []handlers.Link{
			{
				Text: "Basic auth protected (demo:demo)",
				Url:  "/basic",
			},
			{
				Text: "Basic auth protected using DB demoUsers (test@mail.com:1234)",
				Url:  "/basic-auth-db",
			},
			{
				Text: "cookie based protected page",
				Url:  "/cookie",
			},
			{
				Text: "cookie based login (demo:demo)",
				Url:  "/cookie-login",
			},
			{
				Text: "User handling",
				Url:  "/user",
			},
		},
	}
	r.Path("/").Handler(&rootPage)

	return &MyAppHandler{
		router: r,
	}, nil
}

func sampleUserManager() (*user.DbManager, error) {

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return nil, err
	}

	// set some options
	opts := user.ManagerOpts{
		BcryptDifficulty: bcrypt.MinCost,
	}

	userMng, err := user.NewDbManager(db, opts)
	if err != nil {
		return nil, err
	}

	// create a user
	err = userMng.CreateUser(user.User{
		Name:  "test",
		Email: "test@mail.com",
		Pw:    "1234",
	})
	if err != nil {
		return nil, err
	}
	return userMng, nil
}
