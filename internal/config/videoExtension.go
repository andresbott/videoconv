package config

import (
	"errors"
	"fmt"
)

// video Extension
type videoExtension struct {
	extension        string
	videoSettingName []string
}

func (v *videoExtension) Ext() string {
	return v.extension
}

func NewVideoExtension(i interface{}) (*videoExtension, error) {

	item := videoExtension{}
	elementsMap := i.(map[interface{}]interface{})

	for k, v := range elementsMap {
		if k == "extension" {
			item.extension = fmt.Sprintf("%v", v)
			continue
		}

		if k == "videoSettingName" {
			if v != nil {
				settingsMap := v.([]interface{})
				settingsItem := []string{}

				for _, setting := range settingsMap {
					settingsItem = append(settingsItem, fmt.Sprintf("%v", setting))
				}

				item.videoSettingName = settingsItem
				continue
			}
		}
	}

	if item.extension == "" {
		return nil, errors.New("\"extension\" not defined in video setting")
	}

	if len(item.videoSettingName) == 0 {
		return nil, errors.New("\"videoSettingName\" not defined in video setting")
	}

	return &item, nil
}
