package videconv

import (
	"errors"
	"fmt"
	"time"
)

type Profile struct {
	name      string
	scale     string
	duration  time.Duration
	extension string
	codec     string
}

const (
	defCodec     = "h264"
	defExtension = "mp4"
)

func newProfile(in interface{}) (Profile, error) {

	vs := Profile{
		codec:     defCodec,
		extension: defExtension,
	}

	itemsMap := in.(map[interface{}]interface{})
	for k, v := range itemsMap {

		switch k {
		case "name":
			vs.name = fmt.Sprintf("%s", v)
			continue

		case "scale":
			vs.scale = fmt.Sprintf("%s", v)
			continue

		case "duration":

			dur := fmt.Sprintf("%s", v)
			if dur != "" {

			}
			durT, err := time.ParseDuration(dur)
			if err != nil {
				return vs, err
			}
			vs.duration = durT
			continue

		case "extension":
			vs.extension = fmt.Sprintf("%s", v)
			continue

		case "codec":
			vs.codec = fmt.Sprintf("%s", v)
			continue
		}
	}

	if vs.name == "" {
		return vs, errors.New("name not provided for video setting")
	}

	return vs, nil

}
