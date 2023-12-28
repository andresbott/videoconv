package videoconv

import (
	"fmt"
	"github.com/AndresBott/videoconv/app/videoconv/config"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Converter struct {
	Cfg        config.Conf
	DaemonMode bool
	Runner     func(video string)
}

func New(cfg config.Conf) (*Converter, error) {
	c := Converter{
		Cfg: cfg,
	}

	// log level
	lv, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		lv = log.InfoLevel
	}
	log.SetLevel(lv)

	return &c, nil
}

func NewFromCfg(cfgFile string) (*Converter, error) {
	cfg, err := config.NewFromFile(cfgFile)
	if err != nil {
		return nil, err
	}
	return New(cfg)
}

func (vc *Converter) Check(create bool) error {

	for _, item := range vc.Cfg.Locations {
		e := checkLocation(vc.Cfg.ConfigLocation, item, create)
		if e != nil {
			return e
		}
	}
	return nil
}

func checkLocation(cfgLocation string, location config.Location, create bool) error {
	itemPath, err := filepath.Abs(filepath.Join(filepath.Dir(cfgLocation), location.Path))
	if err != nil {
		return fmt.Errorf("error generating absolute path for location: %v", err)
	}
	log.Debug("checking location:" + itemPath)

	if _, err := os.Stat(itemPath); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist or is not accessible", itemPath)
	}

	dirs := []string{
		location.InputDir,
		location.OutputDir,
		location.TmpDir,
		location.FailDir,
	}

	for _, d := range dirs {
		dir, err := filepath.Abs(filepath.Join(itemPath, d))
		if err != nil {
			return fmt.Errorf("error generating absolute path for location: %v", err)
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if !create {
				return fmt.Errorf("directory \"%s\" does not exits", dir)
			}

			log.Warnf("directory %s does not exist or is not accessible, trying to create", dir)
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return fmt.Errorf("unable to create dir %s, %v", dir, err)
			}
		}
	}
	return nil
}

// Run executes the main conversion loop, will exit if not run in daemon mode
func (vc *Converter) Run() {
	log.Info("starting video conversion...")
	for {
		for _, item := range vc.Cfg.Locations {
			if vc.Runner == nil {
				vc.Runner = func(video string) {
					log.Warnf("video processor function not defined, runniung dummy for videon: %s", video)
				}
			}
			vc.runLocation(item, vc.Runner)
		}
		if !vc.DaemonMode {
			log.Info("finished, exiting...")
			break
		}
		log.Infof("finished, sleeping for %s", vc.Cfg.Sleep)
		time.Sleep(vc.Cfg.Sleep)
	}
}

// convert videos on one location
func (vc *Converter) runLocation(location config.Location, callback func(video string)) {
	log.Debug("checking location:" + location.Path)
	locationPath, err := filepath.Abs(filepath.Join(filepath.Dir(vc.Cfg.ConfigLocation), location.Path))
	if err != nil {
		log.Errorf("error generating absolute path for location \"%s\", skipping location, error: %v", locationPath, err)
		return
	}
	if _, err := os.Stat(locationPath); os.IsNotExist(err) {
		log.Errorf("directory %s does not exist or is not accessible, skipping location", locationPath)
		return
	}

	err = checkLocation(vc.Cfg.ConfigLocation, location, false)
	if err != nil {
		log.Errorf("location contains error:\"%s\", sskipping...", err.Error())
	}

	videos, err := findVideos(locationPath, vc.Cfg.VideoExtensions)
	if err != nil {
		log.Errorf("error searching for videos:\"%s\", sskipping...", err.Error())
	}
	for _, video := range videos {
		callback(video)
	}
}

// findVideos recursively searches videos in the rootPath and returns an array of relative paths of videos
func findVideos(rootPath string, videoExtensions []string) ([]string, error) {
	var videos []string

	err := filepath.Walk(rootPath, func(fPath string, fInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if fInfo.IsDir() {
			return nil
		}
		ext := filepath.Ext(fPath)
		ext = ext[1:]
		if !isVideo(ext, videoExtensions) {
			return nil
		}

		rel, err := filepath.Rel(rootPath, fPath)
		if err != nil {
			return err
		}
		videos = append(videos, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return videos, nil
}

// used to check if file is a video based on extension
func isVideo(val string, videoExtensions []string) bool {
	c := strings.TrimSpace(val)
	c = strings.ToLower(val)

	for _, item := range videoExtensions {
		if strings.ToLower(item) == c {
			return true
		}
	}
	return false
}
