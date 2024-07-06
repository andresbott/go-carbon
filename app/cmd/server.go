package cmd

import (
	"git.andresbott.com/Golang/carbon/app/router"
	"git.andresbott.com/Golang/carbon/libs/factory"
	"git.andresbott.com/Golang/carbon/libs/http/handlers"
	"git.andresbott.com/Golang/carbon/libs/http/server"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbFile = "carbon.db"

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
		RunE: func(cmd *cobra.Command, args []string) error {

			// ideally the command reads arguments, loads configuration
			// and creates the server accordingly
			// in this case the command is opinionated

			logOutput, err := factory.ConsoleFileOutput("")
			if err != nil {
				return err
			}
			l := factory.DefaultLogger(factory.InfoLevel, logOutput)

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

			s := server.New(server.Cfg{
				Handler:    rootHandler,
				ObsHandler: handlers.Observability(),
				Logger: func(msg string, isErr bool) {
					if isErr {
						l.Warn().Msg(msg)
					} else {
						l.Info().Msg(msg)
					}
				},
			})

			return s.Start()

		},
	}

	return cmd

}
