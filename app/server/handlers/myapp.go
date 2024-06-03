package handlers

import (
	_ "embed"
	"git.andresbott.com/Golang/carbon/app/spa"
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

	// use session auth
	err := sessionAuthentication(r, demoUsers)
	if err != nil {
		return nil, err
	}

	// user management
	// --------------------------
	// db managed users
	sampleDbUser, err := sampleUserManager()
	if err != nil {
		return nil, err
	}

	userDbHandler, err := userhandler.NewHandler(sampleDbUser)
	if err != nil {
		return nil, err
	}
	r.PathPrefix("/user").Handler(userDbHandler.UserHandler("/user"))

	// root page
	// --------------------------
	demoPage := handlers.SimpleText{
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
	r.Path("/demo").Handler(&demoPage)

	// SPA starts here: ====================================================

	hashKey := []byte("oach9iu2uavahcheephi4FahzaeNge8yeecie4jee9rah9ahrah6tithai7Oow5U")
	blockKey := []byte("eeth3oon5eewifaogeibieShey5eiJ0E")

	sessStor, err := auth.FsStore("", hashKey, blockKey)
	if err != nil {
		return nil, err
	}

	sessionAuth, err := auth.NewSessionMgr(auth.SessionCfg{
		Store: sessStor,
	})

	loginHandler := auth.JsonAuthHandler(sessionAuth, demoUsers)
	r.Path("/login").Methods(http.MethodPost).Handler(loginHandler)

	// load the SPA page
	spaHandler, err := spa.NewCarbonSpa("/")
	if err != nil {
		return nil, err
	}
	_ = spaHandler

	r.Methods(http.MethodGet).PathPrefix("/").Handler(spaHandler)

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
