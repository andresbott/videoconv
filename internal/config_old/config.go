package config_old

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strconv"
)

type Conf struct {
	InputFolder      string
	OutputFolder     string
	TmpDir           string
	IgnoreDir        string
	videoExtensions  videoExtensions
	ProcessorThreads int
	videoSettings    videoSettings
	PollInterval     int
}

// return the video extensions as slice
func (c *Conf) Ext() []string {
	r := []string{}
	for _, i := range c.videoExtensions.Items() {
		r = append(r, i.extension)
	}
	return r
}

// return a list of video settings from the configuration depending on the extension
func (c *Conf) VideoSettingsByExt(ext string) ([]VideoSetting, error) {

	settingNamesFromExtension := c.videoExtensions.GetNamesByExt(ext)
	if settingNamesFromExtension == nil {
		return nil, errors.New("extension: " + ext + " has no video settings defined")
	}

	r := []VideoSetting{}
	for _, name := range settingNamesFromExtension {

		item := c.videoSettings.ItemByName(name)
		if item != nil {
			r = append(r, *item)
		} else {
			log.Warn("video setting: " + name + " does not exist")
		}
	}

	if len(r) == 0 {
		return nil, errors.New("no video settings found for extension: " + ext)
	}

	return r, nil

}

func NewConfig() Conf {

	viper.SetConfigName(configFileName) // name of config file (without extension)
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		log.Fatal("Fatal error config file: %s \n", err)
	}
	c := Conf{}

	// input dir
	in := viper.GetString("input_dir")
	if in == "" {
		log.Fatal("No input dir defined")
	}

	abspath, err := filepath.Abs(in)
	if err != nil {
		log.Fatal("problem with input dir: " + err.Error())
	} else {
		log.Info("Using input dir: " + abspath)
		c.InputFolder = abspath
	}

	// output dir
	out := viper.GetString("output_dir")
	if out == "" {
		log.Fatal("No output dir defined")
	}

	abspath, err = filepath.Abs(out)
	if err != nil {
		log.Fatal("problem with output dir: " + err.Error())
	} else {
		log.Info("Using output dir: " + abspath)
		c.OutputFolder = abspath
	}

	// tmp dir
	tmp := viper.GetString("tmp_dir")
	if out == "" {
		log.Fatal("No tmp dir defined")
	}

	abspath, err = filepath.Abs(tmp)
	if err != nil {
		log.Fatal("problem with tmp dir: " + err.Error())
	} else {
		log.Info("Using tmp dir: " + abspath)
		c.TmpDir = abspath
	}

	// ignore dir
	ignore := viper.GetString("ignore_dir")
	if out == "" {
		log.Fatal("No ignore dir defined")
	}
	abspath, err = filepath.Abs(ignore)
	if err != nil {
		log.Fatal("problem with ignore dir: " + err.Error())
	} else {
		log.Info("Using ignore dir: " + abspath)
		c.IgnoreDir = abspath
	}

	// video in extensions
	vExts := NewVideoExtensions(defaultVideoExtension)
	confVExts := viper.Get("video_extensions")
	if confVExts != nil {
		confVExtsMap := confVExts.([]interface{})

		if len(confVExtsMap) > 0 {
			for _, va := range confVExtsMap {
				vExt, err := NewVideoExtension(va)
				if err != nil {
					log.Warn(err.Error() + " while reading video settings")
					continue
				}
				vExts.addItem(vExt)
			}
		}
	}
	c.videoExtensions = *vExts

	// video settings
	videoSettings := NewVideoSettings(defaultVideoSettings)

	vs := viper.Get("video_settings")
	itemMap := vs.([]interface{})
	for _, v := range itemMap {

		vSetting, err := NewVideoSetting(v)
		if err != nil {
			log.Warn(err.Error() + " while reading video settings")
			continue
		}
		videoSettings.addItem(vSetting)
	}
	c.videoSettings = *videoSettings

	// processor threads
	c.ProcessorThreads = viper.GetInt("threads")
	if c.ProcessorThreads == 0 {
		log.Warn("threads improperly defined, using default: " + string(defaultThread) + " threads")
		c.ProcessorThreads = defaultThread
	}

	// poll interval
	c.PollInterval = viper.GetInt("poll_interval")
	if c.PollInterval == 0 {
		log.Warn("poll interval improperly defined, using default: " + strconv.Itoa(defaultPollInterval) + " seconds")
		c.PollInterval = defaultPollInterval
	} else {
		log.Info("Using poll interval: " + strconv.Itoa(c.PollInterval) + " seconds")
	}

	return c
}
