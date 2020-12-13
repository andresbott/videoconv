package config2

import (
	"errors"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
)

type ConfHandler struct {
	locations    []Location
	profiles     map[string]Profile
	logLevel     string
	threads      int
	pollInterval time.Duration
	ffmpegBin    string
}

const (
	defaultLogLevel     = "info"
	defaultPollDuration = "5m"
	defaultThreads      = 1
	defaultFFmpeg       = "/usr/bin/ffmpeg"
)

// Load uses viper to load the main yaml configuration file
func (cfg *ConfHandler) Load(file string) error {

	fileAbsPath, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigFile(fileAbsPath)

	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	// LogLevel
	cfg.logLevel = v.GetString("log_level")
	if cfg.logLevel == "" {
		cfg.logLevel = defaultLogLevel
	}

	// Poll duration
	poll := v.GetString("poll_interval")
	if poll == "" {
		poll = defaultPollDuration
	}

	cfg.pollInterval, err = time.ParseDuration(poll)
	if err != nil {
		return err
	}

	// ffmpeg
	cfg.ffmpegBin = v.GetString("ffmpeg")
	if cfg.ffmpegBin == "" {
		cfg.ffmpegBin = defaultFFmpeg
	}

	// Threads
	cfg.threads = v.GetInt("threads")
	if cfg.threads == 0 {
		cfg.threads = defaultThreads
	}

	// Load video locations
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
		cfg.locations = append(cfg.locations, *item)
	}

	// load Default Video settings

	cfg.profiles = make(map[string]Profile)

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
		cfg.profiles[item.name] = item
	}

	return nil

}
