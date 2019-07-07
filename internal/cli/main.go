package cli

import (
	"git.andresbott.com/Utilities/go-video-process/internal/config"
	"git.andresbott.com/Utilities/go-video-process/internal/transcode"
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
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

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
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
	Use: "start",
	Short: "Start the daemon",
	Run: func(cmd *cobra.Command, args []string) {


		tr := transcode.NewTranscoder()
		inte := tr.Interval()
		tr.Run()

		for {
			tr.Run()
			time.Sleep(time.Duration(inte) * time.Second)
		}
	},
}

