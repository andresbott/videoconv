package cmd

import (
	"fmt"
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
		RunE: func(cmd *cobra.Command, args []string) error {

			configFile, _ := filepath.Abs("./videoconv.yaml")
			if _, err := os.Stat(configFile); err == nil {
				log.Warn("destination already initialized, skipping...")
				return nil
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
					return fmt.Errorf("unable to create abspath: %s", err)
				}

				log.Infof("creating dir: %s", dir)
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					return fmt.Errorf("unable to create dir %s, %v", dir, err)
				}
			}

			// create templates dir
			dir, err := filepath.Abs("./templates")
			if err != nil {
				return fmt.Errorf("unable to create abspath: %s", err)
			}

			log.Infof("creating dir: %s", dir)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("unable to create dir %s, %v", dir, err)
			}

			tmplFile, _ := filepath.Abs("./templates/empty.tmpl.json")
			tmplContent := config.SampleTmpl()
			log.Infof("writing a sample template file: %s ", tmplFile)
			err = os.WriteFile(tmplFile, []byte(tmplContent), 0644)
			if err != nil {
				return fmt.Errorf("unable to create configuration file %s, %v", tmplFile, err)
			}

			configContent := config.SampleCfg()
			log.Infof("writing configuration file: %s ", configFile)
			err = os.WriteFile(configFile, []byte(configContent), 0644)
			if err != nil {
				return fmt.Errorf("unable to create configuration file %s, %v", configFile, err)
			}
			return nil
		},
	}

	return &cmd

}
