package videconv

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type App struct {
	ConfigFile      string
	logLevel        string
	ffmpegBin       string
	videoExtensions []string
	threads         int
	sleep           time.Duration
	locations       []location
	profiles        map[string]Profile
}

// start the video conversion app
func (vc *App) Start() error {

	err := vc.loadConfig()
	if err != nil {
		return err
	}
	log.Info("starting video conversion")
	vc.loop()
	return nil
}

// prepare the app to run
func (vc *App) loop() {
	for true {
		for _, l := range vc.locations {
			if runLocation(l) {
				continue
			}
		}
		time.Sleep(vc.sleep)
	}
}

// execute one video on one location and return true after a video conversion is finished
func runLocation(l location) bool {
	log.Debug("checking location:" + l.path)

	dir, _ := filepath.Abs(l.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Warnf("directory %s does not exist or is not accessible, skipping location", dir)
		return false
	}

	dirs := []string{
		"./",
		l.inputDir,
		l.outputDir,
		l.tmpDir,
		l.failDir,
	}

	for _, d := range dirs {
		dir, _ := filepath.Abs(l.path + "/" + d)
		r := false
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			r = true
			log.Warnf("directory %s does not exist or is not accessible, skipping location", dir)
		}
		if r {
			return false
		}
	}

	return true
}
