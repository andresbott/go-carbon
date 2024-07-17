package cmd

import (
	"git.andresbott.com/Golang/carbon/app/config"
	"git.andresbott.com/Golang/carbon/app/router"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/factory"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/http/server"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbFile = "carbon.db"

func serverCmd() *cobra.Command {
	var configFile = "./config.yaml"
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := config.Get(configFile)
			if err != nil {
				return err
			}

			// setup the logger
			logOutput, err := factory.ConsoleFileOutput("")
			if err != nil {
				return err
			}
			l := factory.DefaultLogger(factory.GetLogLevel(cfg.Log.Level), logOutput)

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

			// Main APApplication handler
			appCfg := router.AppCfg{
				Logger:      l,
				Db:          db,
				AuthSession: sessionAuth,
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

		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", configFile, "config file")

	return cmd

}
