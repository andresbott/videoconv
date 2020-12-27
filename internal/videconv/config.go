package videconv

import (
	"errors"
	transcoder "github.com/AndresBott/videoconv/internal/transcode"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultLogLevel     = "info"
	defaultPollDuration = "5m"
	DefaultThreads      = 1
	DefaultFFmpeg       = "/usr/bin/ffmpeg"
	defaultOverlayFname = "videoconv.yaml"
)

// prepare the app to run
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

	// load Default Video profiles
	vc.profiles = make(map[string]transcoder.FfmpegOpts)

	confVidSettings := v.Get("profiles")
	if confVidSettings == nil {
		return errors.New("video settings not defined")
	}

	confVidSettList := confVidSettings.([]interface{})
	if len(confVidSettList) == 0 {
		return errors.New("video settings not defined")
	}

	for _, vSet := range confVidSettList {
		item, err := transcoder.NewFromInterface(vSet)
		if err != nil {
			return err
		}
		vc.profiles[item.Name] = *item
	}

	log.Info("loaded config file: " + fileAbsPath)
	log.Infof("using video extensions: %s", strings.Join(vc.videoExtensions, ","))
	return nil
}
