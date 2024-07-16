package cmd

import (
	"git.andresbott.com/Golang/carbon/app/router"
	"git.andresbott.com/Golang/carbon/libs/config"
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

			type Msg struct {
				level string
				msg   string
			}
			configMsg := []Msg{}

			cfg := appCfg{}
			_, err := config.Load(
				config.Defaults{Item: DefaultCfg},
				config.EnvVar{Prefix: "CARBON"},
				config.Unmarshal{Item: &cfg},
				config.Writer{Fn: func(level, msg string) {
					if level == config.InfoLevel {
						configMsg = append(configMsg, Msg{level: "info", msg: msg})
					}
					if level == config.DebugLevel {
						configMsg = append(configMsg, Msg{level: "debug", msg: msg})
					}
				}},
			)
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
			for _, m := range configMsg {
				if m.level == "info" {
					l.Info().Str("component", "config").Msg(m.msg)
				} else {
					l.Debug().Str("component", "config").Msg(m.msg)
				}
			}

			db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
				//Logger: zeroGorm.New(l.ZeroLog, zeroGorm.Cfg{IgnoreRecordNotFoundError: true}),
			})
			if err != nil {
				return err
			}

			rootHandler, err := router.NewAppHandler(l, db)
			if err != nil {
				return err
			}

			s, err := server.New(server.Cfg{
				Addr:       cfg.Main.Addr(),
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
