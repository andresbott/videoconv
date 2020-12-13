package cli

import (
	"fmt"
	"github.com/AndresBott/videoconv/internal/config"
	"github.com/AndresBott/videoconv/internal/transcode"
	"github.com/spf13/cobra"
	"os"
)

func Run() {

	CobraCmd := cobra.Command{
		Use:   "videoconv",
		Short: "CLI utility to send slack messages",
		Run: func(cmd *cobra.Command, args []string) {
			// print help if no command is passed
			if len(args) == 0 {
				err := cmd.Help()

				if err != nil {
					fmt.Println(err)
				}
				os.Exit(0)
			}
		},
	}

	CobraCmd.AddCommand(versionCmd)
	CobraCmd.AddCommand(startCmd)

	if err := CobraCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Version)
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the daemon",
	Run: func(cmd *cobra.Command, args []string) {

		tr := transcode.NewTranscoder()
		for {
			tr.Run()
			tr.Sleep()
		}
	},
}
