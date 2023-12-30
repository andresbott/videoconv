package cmd

import (
	"github.com/AndresBott/videoconv/app/videoconv/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func initCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "init",
		Short: "generate a basic configuration and folder structure",
		Run: func(cmd *cobra.Command, args []string) {

			configFile, _ := filepath.Abs("./videoconv.yaml")
			if _, err := os.Stat(configFile); err == nil {
				log.Warn("destination already initialized, skipping...")
				return
			}

			dirs := []string{
				config.DefaultInputDir,
				config.DefaultOutputDir,
				config.DefaultTmpDir,
				config.DefaultFailDir,
			}

			for _, d := range dirs {

				dir, err := filepath.Abs(filepath.Join("./sample", d))
				if err != nil {
					log.Fatalf("unable to create abspath: %s", err)
				}

				log.Infof("creating dir: %s", dir)
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					log.Fatalf("unable to create dir %s, %v", dir, err)
				}
			}

			configContent := config.SampleCfg()
			d1 := []byte(configContent)

			log.Infof("writing configuration file: %s ", configFile)
			err := os.WriteFile(configFile, d1, 0644)
			if err != nil {
				// todo handle err
				return

			}
		},
	}

	return &cmd

}
