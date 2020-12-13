package config2

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Location struct {
	path            string
	inputDir        string
	outputDir       string
	tmpDir          string
	failDir         string
	appliedProfiles []string
}

const (
	defaultInputDir  = "in"
	defaultOutputDir = "out"
	defaultTmpDir    = "tmp"
	defaultFailDir   = "fail"
)

// loadLocation uses a interface, loaded by viper, that contains the data needed for a video conversion location
func newLocation(in interface{}) (*Location, error) {

	loc := Location{
		path:      "",
		inputDir:  defaultInputDir,
		outputDir: defaultOutputDir,
		tmpDir:    defaultTmpDir,
		failDir:   defaultFailDir,
	}

	itemsMap := in.(map[interface{}]interface{})
	for k, v := range itemsMap {

		switch k {
		case "base_path":
			loc.path = fmt.Sprintf("%s", v)
			continue

		case "input":
			loc.inputDir = fmt.Sprintf("%s", v)
			continue

		case "output":
			loc.outputDir = fmt.Sprintf("%s", v)
			continue

		case "tmp":
			loc.tmpDir = fmt.Sprintf("%s", v)
			continue

		case "fail":
			loc.failDir = fmt.Sprintf("%s", v)
			continue

		case "applied":
			profileList := v.([]interface{})
			if len(profileList) == 0 {
				continue
			}
			for _, i := range profileList {
				loc.appliedProfiles = append(loc.appliedProfiles, fmt.Sprintf("%s", i))
			}
			continue
		}

	}

	if loc.path == "" {
		return nil, errors.New("location base path not provided")
	}
	locAbsPath, err := filepath.Abs(loc.path)
	if err != nil {
		return nil, err
	}
	loc.path = locAbsPath

	return &loc, nil

}

// Overlay checks for an overlay configuration file and if found will load it's content
// Overlay returns a copy of the Location struct in order to be able to reload the file on every check
func (loc *Location) Overlay(fname string) (*Location, error) {
	nl := Location{
		path:            loc.path,
		inputDir:        loc.inputDir,
		outputDir:       loc.outputDir,
		tmpDir:          loc.tmpDir,
		failDir:         loc.failDir,
		appliedProfiles: loc.appliedProfiles,
	}

	if fname == "" {
		return &nl, nil
	}

	p := loc.path + "/" + fname
	if filepath.IsAbs(fname) {
		p = fname
	}

	fileAbsPath, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(fileAbsPath)

	err = v.ReadInConfig()

	// ignore overly file not found
	if err, ok := err.(*os.PathError); ok {
		if err.Err.Error() == "no such file or directory" {
			return &nl, nil
		}
	}

	if err != nil {
		return nil, err
	}

	in := v.GetString("input")
	if in != "" {
		nl.inputDir = in
	}

	out := v.GetString("output")
	if out != "" {
		nl.outputDir = out
	}

	tmp := v.GetString("tmp")
	if tmp != "" {
		nl.tmpDir = tmp
	}

	fail := v.GetString("fail")
	if fail != "" {
		nl.failDir = fail
	}

	drop := v.GetBool("drop_applied")

	if drop {
		nl.appliedProfiles = []string{}
	}

	profileList := v.GetStringSlice("applied")
	nl.appliedProfiles = append(nl.appliedProfiles, profileList...)

	return &nl, nil

}
