package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func Run() {

	cmd := cobra.Command{
		Use:   "videoconv",
		Short: "batch video conversion based on directory observation",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}

	// todo: add global flag debug
	cmd.AddCommand(versionCmd())
	cmd.AddCommand(runCmd())
	cmd.AddCommand(configCmd())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
