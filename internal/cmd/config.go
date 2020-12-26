package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "config",
		Short: "generate configurations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TODO: ")
		},
	}

	return &cmd

}
