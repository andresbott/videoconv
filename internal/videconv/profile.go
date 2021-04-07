package videconv

import "fmt"

const defaultOutExtension = "mp4"

// profile holds the  specific data for every transcoder profile
type profile struct {
	template  string // profile template
	name      string // profile name
	extension string // output extension
}

// newProfile uses the cobra input as interface to generate a profile struct
func newProfile(conf interface{}) (*profile, error) {

	profile := profile{
		extension: defaultOutExtension,
	}
	itemsMap, ok := conf.(map[interface{}]interface{})

	if !ok {
		return nil, fmt.Errorf("unable to cast input from map[interface{}]interface{}")
	}
	for k, v := range itemsMap {
		switch k {
		case "name":
			profile.name = fmt.Sprintf("%s", v)
			continue
		case "extension":
			ext := fmt.Sprintf("%s", v)
			if ext != "" {
				profile.extension = ext
			}
			continue
		case "template":
			profile.template = fmt.Sprintf("%s", v)
			if profile.template == "" {
				return nil, fmt.Errorf("template cannot be empty")
			}
			continue
		}

	}
	return &profile, nil
}
