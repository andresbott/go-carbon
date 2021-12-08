package cmd

import (
	"git.andresbott.com/Golang/carbon/internal/server"
	"git.andresbott.com/Golang/carbon/libs/log/zero"
	"git.andresbott.com/Golang/carbon/libs/log/zeroGorm"
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

			l := zero.NewZero(zero.InfoLevel, nil)

			db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
				Logger: zeroGorm.New(l.ZeroLog, zeroGorm.Cfg{IgnoreRecordNotFoundError: true}),
			})
			if err != nil {
				return err
			}

			s := server.NewServer(server.Cfg{
				Logger: l,
				Db:     db,
			})
			return s.Start()
		},
	}

	return cmd

}
