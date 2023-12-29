package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultLogLevel        = "info"
	defaultSleep           = "5m"
	DefaultFFmpeg          = "/usr/bin/ffmpeg"
	DefaultFFprobe         = "/usr/bin/ffprobe"
	DefaultVideoExtensions = "avi,mkv,mov"
)

type Conf struct {
	// system settings
	LogLevel        string
	Sleep           time.Duration
	FfmpegPath      string
	FfprobePath     string
	VideoExtensions []string
	ConfigLocation  string

	// locations
	Locations []Location
}

func NewFromFile(configFile string) (Conf, error) {
	absCfg, err := filepath.Abs(configFile)
	if err != nil {
		return Conf{}, fmt.Errorf("error calculating absolute path of config file: %v", err)
	}
	c := Conf{
		ConfigLocation: absCfg,
	}
	err = c.loadFile(absCfg)
	return c, err
}

func (cfg *Conf) loadFile(configFile string) error {
	fileAbsPath, err := filepath.Abs(configFile)
	if err != nil {
		return err
	}

	v := viper.New()
	v.SetConfigFile(fileAbsPath)

	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	err = cfg.systemSettings(v)
	if err != nil {
		return err
	}

	// load video locations
	cfg.Locations, err = locationSettings(v)
	if err != nil {
		return err
	}

	return nil
}

// systemSettings loads basic stuff that should already contain sensible defaults
func (cfg *Conf) systemSettings(v *viper.Viper) error {
	// LogLevel
	cfg.LogLevel = v.GetString("log_level")
	if cfg.LogLevel == "" {
		cfg.LogLevel = defaultLogLevel
	}

	// Sleep duration in daemon mode
	poll := v.GetString("poll_interval")
	if poll == "" {
		poll = defaultSleep
	}
	sleep, err := time.ParseDuration(poll)
	if err != nil {
		return err
	}
	cfg.Sleep = sleep

	// ffmpeg
	cfg.FfmpegPath = v.GetString("ffmpeg")
	if cfg.FfmpegPath == "" {
		cfg.FfmpegPath = DefaultFFmpeg
	}
	if _, err := os.Stat(cfg.FfmpegPath); os.IsNotExist(err) {
		return fmt.Errorf("ffmpeg not found on Path: %s", cfg.FfmpegPath)
	}

	// ffprobe
	cfg.FfprobePath = v.GetString("ffprobe")
	if cfg.FfprobePath == "" {
		cfg.FfprobePath = DefaultFFprobe
	}
	if _, err := os.Stat(cfg.FfprobePath); os.IsNotExist(err) {
		return fmt.Errorf("ffprobe not found on Path: %s", cfg.FfprobePath)
	}

	// Video Extensions
	cfg.VideoExtensions = v.GetStringSlice("video_extensions")
	if len(cfg.VideoExtensions) == 0 {
		cfg.VideoExtensions = strings.Split(DefaultVideoExtensions, ",")
	}

	return nil
}

type Location struct {
	Path      string
	InputDir  string
	OutputDir string
	TmpDir    string
	FailDir   string
	Profiles  []Profile
}

const (
	DefaultInputDir  = "in"
	DefaultOutputDir = "out"
	DefaultTmpDir    = "tmp"
	DefaultFailDir   = "fail"
)

func locationSettings(viper *viper.Viper) ([]Location, error) {
	confLocations := viper.Get("locations")
	if confLocations == nil {
		return nil, errors.New("video locations not defined")
	}

	confLocList := confLocations.([]interface{})
	if len(confLocList) == 0 {
		return nil, errors.New("video locations does not contain any entry")
	}
	locations := []Location{}

	for _, item := range confLocList {
		loc, err := buildLocation(item)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil

}

func buildLocation(in interface{}) (Location, error) {
	loc := Location{
		Path:      "",
		InputDir:  DefaultInputDir,
		OutputDir: DefaultOutputDir,
		TmpDir:    DefaultTmpDir,
		FailDir:   DefaultFailDir,
	}

	for k, v := range in.(map[interface{}]interface{}) {

		switch k {
		case "path":
			loc.Path = fmt.Sprintf("%s", v)
			continue

		case "input":
			loc.InputDir = fmt.Sprintf("%s", v)
			continue

		case "output":
			loc.OutputDir = fmt.Sprintf("%s", v)
			continue

		case "tmp":
			loc.TmpDir = fmt.Sprintf("%s", v)
			continue

		case "fail":
			loc.FailDir = fmt.Sprintf("%s", v)
			continue

		case "profiles":
			profileList := v.([]interface{})
			if len(profileList) == 0 {
				continue
			}
			for _, i := range profileList {

				got, err := buildProfile(i)
				if err != nil {
					return Location{}, err
				}

				loc.Profiles = append(loc.Profiles, got)
			}
			continue
		}
	}

	// verify the Location is configured properly

	if loc.Path == "" {
		return loc, errors.New("location path not provided")
	}

	return loc, nil
}

type Profile struct {
	Template string
	Args     map[string]string
}

func buildProfile(in interface{}) (Profile, error) {
	pr := Profile{
		Args: map[string]string{},
	}
	// iterate over the profile entry
	for k, v := range in.(map[interface{}]interface{}) {

		// stringify the value
		value := ""
		switch v.(type) {
		case int:
			value = fmt.Sprintf("%d", v)
		case string:
			value = v.(string)
		}

		// parse the keys
		switch k {
		case "template":
			pr.Template = value
			continue
		default:
			pr.Args[k.(string)] = value
		}

	}

	return pr, nil
}

func SampleCfg() string {

	return `# sample configuration file for videconv
log_level: "info"

# poll interval looking for new videos if running in daemon mode
poll_interval: "5m"

# if ffmpeg is in a differnt location
ffmpeg:  "/usr/bin/ffmpeg"
ffprobe: "/usr/bin/ffprobe"

# only files with this extensions are processed
video_extensions:
  - avi
  - mkv
  - mov
  - wmv
  - mp4

# list of locations where to perform video conversions
locations:
  - path: "./sample"
    input:  "in" 	# input directory to put videos
    output: "out"   # output where processed videos are moved
    tmp:    "tmp"   # temporary directory while processing a video
    fail:   "fail"  # failed videos are moved here
    profiles:
      - template: "sample"
        key: "value"

`

}
