package cmd

import (
	"github.com/AndresBott/videoconv/app/videoconv"
	"github.com/AndresBott/videoconv/app/videoconv/config"
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

				cfg, err := config.NewFromFile(configfile)
				if err != nil {
					return err
				}

				vidConv, err := videoconv.New(cfg)
				if err != nil {
					return err
				}
				vidConv.Run()

			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&configfile, "config", "c", configfile, "configuration file")
	cmd.Flags().BoolVarP(&daemon, "daemon", "d", daemon, "run in daemon mode")

	return &cmd
}
