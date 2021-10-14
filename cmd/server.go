package cmd

import "github.com/spf13/cobra"

func serverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
	}

	return cmd

}
