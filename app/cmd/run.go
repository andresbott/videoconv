package cmd

import (
	"github.com/AndresBott/videoconv/internal/videconv"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {

	configfile := "videoconv.yaml"
	daemon := false

	cmd := cobra.Command{
		Use:   "run",
		Short: "run transcoding on the configured destinations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				// main execution path if no args are passed
				app := videconv.App{
					ConfigFile: configfile,
					DaemonMode: daemon,
				}
				return app.Start()
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&configfile, "config", "c", configfile, "configuration file")
	cmd.Flags().BoolVarP(&daemon, "daemon", "d", daemon, "run in daemon mode")

	return &cmd
}
