package cmd

import (
	"git.andresbott.com/Golang/carbon/internal/server"
	"git.andresbott.com/Golang/carbon/libs/log"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbFile = "carbon.db"

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
		RunE: func(cmd *cobra.Command, args []string) error {

			// ideally the command reads arguments, loads configuration
			// and creates the server accordingly
			// in this case the command is opinionated

			l, err := log.NewZapper()
			if err != nil {
				return err
			}

			// todo gorm logger
			db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})

			s := server.NewServer(server.Cfg{
				Logger: l,
				Db:     db,
			})
			return s.Start()
		},
	}

	return cmd

}
