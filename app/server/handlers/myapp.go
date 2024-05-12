package handlers

import (
	_ "embed"
	"git.andresbott.com/Golang/carbon/internal/http/userhandler"
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

	// static demos users
	demoUsers := user.StaticUsers{
		Users: map[string]string{
			"demo": "demo",
		},
	}
	// add basic auth with fixed users
	basicAuth(r, demoUsers)

	// use session auth
	err := sessionAuthentication(r, demoUsers)
	if err != nil {
		return nil, err
	}

	// db managed users
	sampleDbUser, err := sampleUserManager()
	if err != nil {
		return nil, err
	}
	// add basic auth with users from a in-memory DB
	basicAuthDb(r, sampleDbUser)

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
				Text: "session based protected page",
				Url:  sessionContent,
			},
			{
				Text: "session based login (demo:demo)",
				Url:  sessionLogin,
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
