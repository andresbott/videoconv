package videconv

import (
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
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
	OverlayFname    string
	profiles        map[string]Profile
}

// start the video conversion app
func (vc *App) Start() error {

	err := vc.loadConfig()
	if err != nil {
		return err
	}
	log.Info("starting video conversion...")
	vc.loop()
	return nil
}

// prepare the app to run
func (vc *App) loop() {
	for true {
		for _, l := range vc.locations {
			if vc.runLocation(l) {
				continue
			}
		}
		time.Sleep(vc.sleep)
	}
}

// execute one video on one location and return true after a video conversion is finished
func (vc *App) runLocation(l location) bool {
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
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Warnf("directory %s does not exist or is not accessible, trying to create", dir)
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return false
			}
		}
	}

	newLoc, err := l.loadOverlay(vc.OverlayFname)
	if err != nil {
		log.Errorf("error while loading overlay: %s", err)
		return false
	}
	videos, err := vc.findVideoFiles(newLoc.path + "/" + newLoc.inputDir)
	if err != nil {
		log.Errorf("error while searching videos: %s", err)
		return false
	}
	log.Debugf("found %d, videos: \"%s\"", len(videos), strings.Join(videos, "\", \""))

	//Next, adapt and use transcoder to process the video list
	spew.Dump(videos)

	return true
}

// findVideoFiles recurses the provided root dir and searches video files
// returns a slice of relative paths
func (vc *App) findVideoFiles(rootPath string) ([]string, error) {

	var matches []string

	err := filepath.Walk(rootPath, func(fPath string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}
		ext := filepath.Ext(fPath)
		ext = ext[1:]

		if !vc.isVideo(ext) {
			return nil
		}

		rel, err := filepath.Rel(rootPath, fPath)
		if err != nil {
			return err
		}
		matches = append(matches, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// used to check if file is a video based on extansion
func (vc *App) isVideo(val string) bool {
	c := strings.TrimSpace(val)
	c = strings.ToLower(val)

	for _, item := range vc.videoExtensions {
		if strings.ToLower(item) == c {
			return true
		}
	}
	return false
}
