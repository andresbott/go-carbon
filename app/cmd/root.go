package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// Execute is the entry point for the command line
func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "carbon",
		Short: "carbon is a demo application to explore the framework",
		Long:  "carbon is a demo application to explore the framework",
	}

	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		_ = cmd.Help()
		return nil
	})

	cmd.AddCommand(
		serverCmd(),
	)

	return cmd
}
