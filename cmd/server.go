package cmd

import (
	"git.andresbott.com/Golang/carbon/internal/server"
	"git.andresbott.com/Golang/carbon/libs/log"
	"github.com/spf13/cobra"
)

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
		RunE: func(cmd *cobra.Command, args []string) error {

			l, err := log.NewZapper()
			if err != nil {
				return err
			}

			s := server.NewServer("", l)
			return s.Start()
		},
	}

	return cmd

}
