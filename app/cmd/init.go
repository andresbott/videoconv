package cmd

import (
	"github.com/AndresBott/videoconv/internal/videconv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func configCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "init",
		Short: "generate a basic configuration and folder structure",
		Run: func(cmd *cobra.Command, args []string) {

			locationName := "main"
			configFile, _ := filepath.Abs("./videoconv.yaml")
			if _, err := os.Stat(configFile); err == nil {
				log.Warn("destination already initialized, skipping...")
				return
			}

			dirs := []string{
				videconv.DefaultInputDir,
				videconv.DefaultOutputDir,
				videconv.DefaultTmpDir,
				videconv.DefaultFailDir,
			}

			for _, d := range dirs {

				dir, err := filepath.Abs(filepath.Join("./"+locationName, d))
				if err != nil {
					log.Fatalf("unable to create abspath: %s", err)
				}

				log.Infof("creating dir: %s", dir)
				err = os.MkdirAll(dir, 0755)
				if err != nil {
					log.Fatalf("unable to create dir %s, %v", dir, err)
				}
			}

			configContent := `
---
# sample configuration
log_level: "info"

# poll interval looking for new videos if running in daemon mode
poll_interval: "30m"

# Location for ffmpeg use /usr/bin/ffmpeg as default
ffmpeg: "` + videconv.DefaultFFmpeg + `"
# amount of threads configured for ffmpeg
threads: ` + strconv.Itoa(videconv.DefaultThreads) + `

# list of video file extensions to handle
video_extensions:
  - mp4
  - wmv
  - mkv
  - avi

# list of locations where to perform video conversions
# see Readme for details
locations:
  - base_path: "./` + locationName + `"
    applied:
      - "720p_sample"

# list of video conversion profiles to be used in the locations
# see Readme for details
profiles:

  - name: 720p_sample
    threads: "1"
    codec: "libx264"
    quality_crf: "23"
    quality_preset: "medium"
    quality_tune: "film"
    scale: 720
    duration: "30s"
    start: "10s"
    cuda_decoding: "false" # set to "true" to use -hwaccel cuda
    cuda_hw_output: "false" # set to true to use -hwaccel_output_format cuda for late nvec encoding
    extra: "" # add extra ffmpeg parameters not listed above
`
			d1 := []byte(configContent)

			log.Infof("writing configuration file: %s ", configFile)
			err := ioutil.WriteFile(configFile, d1, 0644)
			if err != nil {

			}

		},
	}

	return &cmd

}
