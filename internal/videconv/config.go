package videconv

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultLogLevel     = "info"
	defaultPollDuration = "5m"
	DefaultThreads      = 1
	DefaultFFmpeg       = "/usr/bin/ffmpeg"
	DefaultFFprobe      = "/usr/bin/ffprobe"
	defaultOverlayFname = "videoconv.yaml"
)

// loadConfig loads the selected location for the config file and loads the values into the app
func (vc *App) loadConfig() error {

	fileAbsPath, err := filepath.Abs(vc.ConfigFile)
	if err != nil {
		return err
	}

	if vc.OverlayFname == "" {
		vc.OverlayFname = defaultOverlayFname
	}

	v := viper.New()
	v.SetConfigFile(fileAbsPath)

	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	// LogLevel
	vc.logLevel = v.GetString("log_level")
	if vc.logLevel == "" {
		vc.logLevel = defaultLogLevel
	}
	// log level
	lv, err := log.ParseLevel(vc.logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(lv)

	// Poll duration
	poll := v.GetString("poll_interval")
	if poll == "" {
		poll = defaultPollDuration
	}

	vc.sleep, err = time.ParseDuration(poll)
	if err != nil {
		return err
	}

	// ffmpeg
	vc.ffmpegBin = v.GetString("ffmpeg")
	if vc.ffmpegBin == "" {
		vc.ffmpegBin = DefaultFFmpeg
	}
	if _, err := os.Stat(vc.ffmpegBin); os.IsNotExist(err) {
		return fmt.Errorf("ffmpeg not found on path: %s", vc.ffmpegBin)
	}

	// ffprobe
	vc.ffProbeBin = v.GetString("ffprobe")
	if vc.ffProbeBin == "" {
		vc.ffProbeBin = DefaultFFprobe
	}
	if _, err := os.Stat(vc.ffProbeBin); os.IsNotExist(err) {
		return fmt.Errorf("ffprobe not found on path: %s", vc.ffProbeBin)
	}

	// Threads
	vc.threads = v.GetInt("threads")
	if vc.threads == 0 {
		vc.threads = DefaultThreads
	}

	// Video Extensions
	vc.videoExtensions = v.GetStringSlice("video_extensions")

	// load video locations
	confLocations := v.Get("locations")
	if confLocations == nil {
		return errors.New("video locations not defined")
	}

	confLocList := confLocations.([]interface{})
	if len(confLocList) == 0 {
		return errors.New("video locations not defined")
	}

	for _, loc := range confLocList {
		item, err := newLocation(loc)
		if err != nil {
			return err
		}
		vc.locations = append(vc.locations, *item)
	}

	// load Video profiles
	vc.profiles = make(map[string]profile)

	confVidSettings := v.Get("profiles")
	if confVidSettings == nil {
		return errors.New("video settings not defined")
	}

	confVidSettList := confVidSettings.([]interface{})
	if len(confVidSettList) == 0 {
		return errors.New("video settings not defined")
	}

	for _, vSet := range confVidSettList {
		item, err := newProfile(vSet)
		if err != nil {
			return err
		}
		vc.profiles[item.name] = *item
	}

	log.Info("loaded config file: " + fileAbsPath)
	log.Infof("using video extensions: %s", strings.Join(vc.videoExtensions, ","))
	return nil
}
