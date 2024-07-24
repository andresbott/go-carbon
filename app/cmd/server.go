package cmd

import (
	"fmt"
	"git.andresbott.com/Golang/carbon/app/config"
	"git.andresbott.com/Golang/carbon/app/router"
	"git.andresbott.com/Golang/carbon/internal/model/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/http/server"
	"git.andresbott.com/Golang/carbon/libs/logzero"
	"git.andresbott.com/Golang/carbon/libs/user"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

const dbFile = "carbon.db"

func serverCmd() *cobra.Command {
	var configFile = "./config.yaml"
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(configFile)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", configFile, "config file")
	return cmd
}

func runServer(configFile string) error {

	cfg, err := config.Get(configFile)
	if err != nil {
		return err
	}

	// setup the logger
	logOutput, err := logzero.ConsoleFileOutput("")
	if err != nil {
		return err
	}
	l := logzero.DefaultLogger(logzero.GetLogLevel(cfg.Log.Level), logOutput)

	l.Info().Str("version", Version).Str("component", "startup").
		Msgf("running version %s, build date: %s, commint: %s ", Version, BuildTime, ShaVer)

	// print config messages delayed
	for _, m := range cfg.Msgs {
		if m.Level == "info" {
			l.Info().Str("component", "config").Msg(m.Msg)
		} else {
			l.Debug().Str("component", "config").Msg(m.Msg)
		}
	}

	// initialize DB
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		//Logger: zeroGorm.New(l.ZeroLog, zeroGorm.Cfg{IgnoreRecordNotFoundError: true}),
	})
	if err != nil {
		return err
	}

	// session based auth
	//cookieStore, err := auth.CookieStore(hashKey, blockKey)
	cookieStore, err := auth.FsStore(cfg.Auth.SessionPath, []byte(cfg.Auth.HashKey), []byte(cfg.Auth.BlockKey))
	if err != nil {
		return err
	}
	sessionAuth, err := auth.NewSessionMgr(auth.SessionCfg{
		Store: cookieStore,
	})
	if err != nil {
		return err
	}

	var users auth.UserLogin
	// load the correct user manager
	switch cfg.Auth.UserStore.StoreType {
	case "static":
		staticUsers := user.StaticUsers{}
		for _, u := range cfg.Auth.UserStore.Users {
			staticUsers.Add(u.Name, u.Pw)
		}
		users = &staticUsers
		l.Debug().Str("component", "users").Msgf("loading %d static user(s)", len(staticUsers.Users))
	case "file":
		if cfg.Auth.UserStore.FilePath == "" {
			return fmt.Errorf("no path for users file is empty")
		}
		staticUsers, err := user.FromFile(cfg.Auth.UserStore.FilePath)
		if err != nil {
			return err
		}
		users = staticUsers
		l.Debug().Str("component", "users").Msgf("loading %d users from file", len(staticUsers.Users))
	default:
		return fmt.Errorf("wrong user store in configuration, %s is not supported", cfg.Auth.UserStore.StoreType)
	}

	// init task manager
	taskMngr, err := tasks.New(db, &sync.Mutex{})
	if err != nil {
		return fmt.Errorf("unable to create task manager :%v", err)
	}
	// Main APApplication handler
	appCfg := router.AppCfg{
		Logger:   l,
		Db:       db,
		AuthMngr: sessionAuth,
		Users:    users,
		Tasks:    taskMngr,
	}
	rootHandler, err := router.NewAppHandler(appCfg)
	if err != nil {
		return err
	}

	s, err := server.New(server.Cfg{
		Addr:       cfg.Server.Addr(),
		Handler:    rootHandler,
		SkipObs:    false,
		ObsAddr:    cfg.Obs.Addr(),
		ObsHandler: handlers.Observability(),
		Logger: func(msg string, isErr bool) {
			if isErr {
				l.Warn().Str("component", "server").Msg(msg)
			} else {
				l.Info().Str("component", "server").Msg(msg)
			}
		},
	})
	if err != nil {
		return err
	}

	return s.Start()
}
