package cmd

import (
	"fmt"
	"github.com/AndresBott/videoconv/internal/videconv"
	"github.com/spf13/cobra"
	"os"
)

const Version = "0.2-SNAPSHOT"

func Run() {

	configfile := "videoconv.yaml"

	cmd := cobra.Command{
		Use:   "videoconv",
		Short: "batch video conversion based on directory observation",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) == 0 {
				// main execution path if no args are passed
				app := videconv.App{
					ConfigFile: configfile,
				}
				return app.Start()
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&configfile, "config", "c", configfile, "configuration file")

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	})

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
