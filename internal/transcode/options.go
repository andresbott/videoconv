package transcoder

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var allowedPresets = []string{
	"ultrafast",
	"superfast",
	"veryfast",
	"faster",
	"fast",
	"medium",
	"slow",
	"slower",
	"veryslow",
}
var allowedTune = []string{
	"film",
	"grain",
	"animation",
	"hq",       // h264_nvenc
	"ll",       // h264_nvenc
	"ull",      // h264_nvenc
	"lossless", // h264_nvenc
}

var allowedVideoCodec = []string{
	"libx264",
	"h264_nvenc",
	"hevc_nvenc",
}
var defaultExtension = "mp4"

type FfmpegOpts struct {
	Threads       int    `attr:"-threads"`
	VideoCodec    string `attr:"-c:v"`
	QualityCRF    *int   `attr:"-crf"`
	QualityPreset string `attr:"-preset"`
	QualityTune   string `attr:"-tune"`
	VideoScale    int
	VideoDuration time.Duration `attr:"-t"`
	VideoStart    time.Duration `attr:"-ss"`
	CudaDecoding  bool
	CudaHwOutput  bool

	Extra string `attr:" "`

	Name           string
	videoExtension string

	once sync.Once
}

const (
	NameNotProvidedError = "name not provided for ffmpeg options"
	NotAllowedCodec      = "codec is not in the allowed list"
	NotAllowedPreset     = "video preset is not in the allowed list"
	NotAllowedTune       = "video tune is not in the allowed list"
)

// NewFromInterface takes an interface, generated from a yaml loader, to generate an ffmpegOpts struct file
func NewFromInterface(in interface{}) (*FfmpegOpts, error) {
	opts := FfmpegOpts{}

	itemsMap := in.(map[interface{}]interface{})
	for k, v := range itemsMap {

		switch k {
		case "name":
			opts.Name = fmt.Sprintf("%s", v)
			continue

		case "threads":
			i, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("wrong value for threads: %v", err))
			}
			opts.Threads = int(i)
			continue

		case "codec":
			codec := fmt.Sprintf("%s", v)
			if !inSlice(allowedVideoCodec, codec) {
				return nil, errors.New(NotAllowedCodec)
			}
			opts.VideoCodec = codec
			continue

		case "quality_crf":
			i, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("wrong value for quality: %v", err))
			}
			val := int(i)
			opts.QualityCRF = &val
			// verify crf value
			if opts.QualityCRF != nil && *opts.QualityCRF > 51 {
				max := 51
				opts.QualityCRF = &max
			} else if opts.QualityCRF != nil && *opts.QualityCRF < 0 {
				min := 0
				opts.QualityCRF = &min
			}
			continue

		case "quality_preset":
			val := fmt.Sprintf("%s", v)
			if !inSlice(allowedPresets, val) {
				return nil, errors.New(NotAllowedPreset)
			}
			opts.QualityPreset = val
			continue

		case "quality_tune":
			val := fmt.Sprintf("%s", v)
			if !inSlice(allowedTune, val) {
				return nil, errors.New(NotAllowedTune)
			}
			opts.QualityTune = val
			continue

		case "scale":
			i, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("wrong value for scale: %v", err))
			}
			opts.VideoScale = int(i)

			continue

		case "duration":

			d, err := time.ParseDuration(fmt.Sprintf("%v", v))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("wrong value for duratiion: %v", err))
			}
			opts.VideoDuration = d

			continue

		case "start":

			d, err := time.ParseDuration(fmt.Sprintf("%v", v))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("wrong value for start: %v", err))
			}
			opts.VideoStart = d

			continue

		case "cuda_decoding":
			val := fmt.Sprintf("%s", v)
			if val == "true" {
				opts.CudaDecoding = true
			}
			continue

		case "cuda_hw_output":
			val := fmt.Sprintf("%s", v)
			if val == "true" {
				opts.CudaHwOutput = true
			}
			continue

		case "extra":
			val := fmt.Sprintf("%s", v)
			opts.Extra = val
			continue
		}

	}

	if opts.Name == "" {
		return &opts, errors.New(NameNotProvidedError)
	}

	return &opts, nil
}

func (opts *FfmpegOpts) VideoExt() string {

	if opts.videoExtension != "" {
		return opts.videoExtension
	}
	return defaultExtension
}

func (opts *FfmpegOpts) Args() ([]string, []string, error) {

	// verify crf value
	if opts.QualityCRF != nil && *opts.QualityCRF > 51 {
		max := 51
		opts.QualityCRF = &max
	} else if opts.QualityCRF != nil && *opts.QualityCRF < 0 {
		min := 0
		opts.QualityCRF = &min
	}

	// verify Quality presets
	if opts.QualityPreset != "" && !inSlice(allowedPresets, opts.QualityPreset) {
		return nil, nil, errors.New("provided preset is not allowed")
	}

	// verify Quality tune
	if opts.QualityTune != "" && !inSlice(allowedTune, opts.QualityTune) {
		return nil, nil, errors.New("provided quality tune is not allowed")
	}

	t := reflect.ValueOf(opts)
	optsItem := reflect.Indirect(t).Interface()

	f := reflect.TypeOf(optsItem)
	v := reflect.ValueOf(optsItem)

	var values []string

	// video scale
	if opts.VideoScale != 0 {
		//`-vf 'scale=-2:min(480\,ih-mod(ih\,2))'
		values = append(values, "-vf", fmt.Sprintf("scale=-2:min(%d\\,ih-mod(ih\\,2))", opts.VideoScale))
	}

	for i := 0; i < f.NumField(); i++ {

		// skip unexported fields
		if !v.Field(i).CanInterface() {
			continue
		}

		attr := f.Field(i).Tag.Get("attr")
		if attr == "" {
			continue
		}

		value := v.Field(i).Interface()
		if v.Field(i).Kind() == reflect.Ptr {
			if !v.Field(i).IsNil() {
				if v, ok := value.(*int); ok {
					values = append(values, attr, fmt.Sprintf("%d", *v))
				}
			}
		} else {
			if !v.Field(i).IsZero() {

				if v, ok := value.(string); ok {
					if strings.TrimSpace(attr) != "" {
						values = append(values, attr)
					}
					splitValues := strings.Split(v, " ")
					for _, spv := range splitValues {
						values = append(values, strings.TrimSpace(spv))
					}
				}

				if _, ok := value.(bool); ok {
					values = append(values, attr)
				}

				if v, ok := value.(int); ok {
					values = append(values, attr, fmt.Sprintf("%d", v))
				}

				if v, ok := value.(time.Duration); ok {
					values = append(values, attr, fmtDuration(v))
				}
			}
		}
	}

	// handle cuda pre args
	var preValues []string
	if opts.CudaDecoding == true {
		preValues = append(preValues, "-hwaccel", "cuda")
	}
	if opts.CudaHwOutput == true {
		preValues = append(preValues, "-hwaccel_output_format", "cuda")
	}

	return dropEmpty(preValues), dropEmpty(values), nil
}

// remove empty items in slice
func dropEmpty(in []string) []string {
	var out []string
	for _, v := range in {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}

func fmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second
	d -= s * time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
