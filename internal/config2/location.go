package config2

import (
	"errors"
	"fmt"
	"path/filepath"
)

type Location struct {
	path      string
	inputDir  string
	outputDir string
	tmpDir    string
	failDir   string
	profiles  []string
}

const (
	defaultInputDir  = "in"
	defaultOutputDir = "out"
	defaultTmpDir    = "tmp"
	defaultFailDir   = "fail"
)

// loadLocation uses a interface, loaded by viper, that contains the data needed for a video conversion location
func newLocation(in interface{}) (Location, error) {

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

		case "profiles":
			profileList := v.([]interface{})
			if len(profileList) == 0 {
				continue
			}
			for _, i := range profileList {
				loc.profiles = append(loc.profiles, fmt.Sprintf("%s", i))
			}
			continue
		}

	}

	if loc.path == "" {
		return loc, errors.New("location base path not provided")
	}
	locAbsPath, err := filepath.Abs(loc.path)
	if err != nil {
		return loc, err
	}
	loc.path = locAbsPath

	return loc, nil

}
