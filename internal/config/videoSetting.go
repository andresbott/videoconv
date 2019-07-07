package config

import (
	"errors"
	"fmt"
)

type VideoSetting struct {
	name         string
	cmd          string
	outExtension string
}

func NewVideoSetting(i interface{}) (*VideoSetting, error) {

	item := VideoSetting{}

	elementsMap := i.(map[interface{}]interface{})

	for k, v := range elementsMap {
		if k == "name" {
			item.name = fmt.Sprintf("%v", v)
			continue
		}

		if k == "cmd" {
			item.cmd = fmt.Sprintf("%v", v)
			continue
		}

		if k == "out_extension" {
			item.outExtension = fmt.Sprintf("%v", v)
			continue
		}
	}

	if item.name == "" {
		return nil, errors.New("\"name\" not defined in video setting")
	}

	if item.cmd == "" {
		return nil, errors.New("\"cmd\" not defined in video setting")
	}

	if item.outExtension == "" {
		return nil, errors.New("\"out_extension\" not defined in video setting")
	}

	return &item, nil
}

func (v *VideoSetting) Name() string {
	return v.name
}

func (v *VideoSetting) Cmd() string {
	return v.cmd
}

func (v *VideoSetting) OutputExtension() string {
	return v.outExtension
}
