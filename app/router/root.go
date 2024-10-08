package router

import (
	_ "embed"
	"github.com/andresbott/go-carbon/app/spa"
	"github.com/andresbott/go-carbon/internal/model/tasks"
	"github.com/andresbott/go-carbon/libs/auth"
	"github.com/andresbott/go-carbon/libs/http/middleware"
	"github.com/andresbott/go-carbon/libs/user"
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

//  go: embed ../handlers/tmpl/loginForm.html
//var loginForm string

type AppCfg struct {
	Logger   *zerolog.Logger
	Db       *gorm.DB
	AuthMngr *auth.SessionMgr
	Users    auth.UserLogin
	Tasks    *tasks.Manager
}

// NewAppHandler generates the main url router handler to be used in the server
func NewAppHandler(cfg AppCfg) (*MyAppHandler, error) {

	r := mux.NewRouter()

	// add observability
	hist := middleware.NewHistogram("", nil, nil)
	r.Use(func(handler http.Handler) http.Handler {
		return middleware.PromLogMiddleware(handler, hist, cfg.Logger)
	})

	// TODO not to have in production
	throttle := middleware.ReqThrottle{
		MinDelay: 1500 * time.Millisecond,
		MaxDelay: 3000 * time.Millisecond,
		On:       true,
	}
	r.Use(throttle.Throttle)

	// attach API v0 handlers
	err := apiV0(r, cfg.AuthMngr, cfg.Users, cfg.Tasks) // api/v0 routes
	if err != nil {
		return nil, err
	}

	// attach the basic auth handler
	err = basicAuthProtected(r.PathPrefix("/basic").Subrouter(), cfg.Users) // api/v0 routes
	if err != nil {
		return nil, err
	}

	// attach spa handler
	// if you want to serve the spa from the root, pass "/" to the spa handler and the path prefix
	// not that the SPA base and route needs to be adjusted accordingly
	spaHandler, err := spa.NewCarbonSpa("/spa")
	if err != nil {
		return nil, err
	}
	r.Methods(http.MethodGet).PathPrefix("/spa").Handler(spaHandler)

	// attach the demo handler on the root path
	err = demo(r)
	if err != nil {
		return nil, err
	}

	// use session auth
	err = SessionProtected(r, cfg.AuthMngr)
	if err != nil {
		return nil, err
	}

	// root page
	// --------------------------

	// SPA starts here: ====================================================

	//hashKey := []byte("oach9iu2uavahcheephi4FahzaeNge8yeecie4jee9rah9ahrah6tithai7Oow5U")
	//blockKey := []byte("eeth3oon5eewifaogeibieShey5eiJ0E")
	//
	//sessStor, err := auth.FsStore("", hashKey, blockKey)
	//if err != nil {
	//	return nil, err
	//}
	//
	//sessionAuth, err := auth.NewSessionMgr(auth.SessionCfg{
	//	Store: sessStor,
	//})
	//
	//loginHandler := auth.JsonAuthHandler(sessionAuth, demoUsers)
	//r.Path("/login").Methods(http.MethodPost).Handler(loginHandler)

	return &MyAppHandler{
		router: r,
	}, nil
}

// TODO put on propper place and usage
// nolint: unused
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
