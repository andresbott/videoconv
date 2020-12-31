package transcoder

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

func getIntPointer(i int) *int {
	a := i
	return &a
}

func TestFfmpegOpst_Args(t *testing.T) {

	tcs := []struct {
		name              string
		in                FfmpegOpts
		expectedInputArgs []string
		expectedArgs      []string
		expectedErr       string
	}{
		{
			name: "Threads", in: FfmpegOpts{Threads: 2},
			expectedArgs: []string{"-threads", "2"},
		},
		{
			name: "VideoCodec", in: FfmpegOpts{VideoCodec: "libx264"},
			expectedArgs: []string{"-c:v", "libx264"},
		},
		{
			name: "VideoQuality", in: FfmpegOpts{QualityCRF: getIntPointer(23)},
			expectedArgs: []string{"-crf", "23"},
		},
		{
			name: "VideoQuality out of bounds", in: FfmpegOpts{QualityCRF: getIntPointer(100)},
			expectedArgs: []string{"-crf", "51"},
		},
		{
			name: "VideoPreset", in: FfmpegOpts{QualityPreset: "medium"},
			expectedArgs: []string{"-preset", "medium"},
		},
		{
			name: "QualityTune", in: FfmpegOpts{QualityTune: "animation"},
			expectedArgs: []string{"-tune", "animation"},
		},
		{
			name: "VideoScale", in: FfmpegOpts{VideoScale: 720},
			expectedArgs: []string{"-vf", `scale=-2:min(720\,ih-mod(ih\,2))`},
		},
		{
			name: "Duration", in: FfmpegOpts{VideoDuration: 2 * time.Minute, VideoStart: 10 * time.Second},
			expectedArgs: []string{"-t", "00:02:00", "-ss", "00:00:10"},
		},
		{
			name: "Cuda", in: FfmpegOpts{CudaDecoding: true, CudaHwOutput: true},
			expectedInputArgs: []string{"-hwaccel", "cuda", "-hwaccel_output_format", "cuda"},
		},
		{
			name: "Extra", in: FfmpegOpts{Extra: "-r 30"},
			expectedArgs: []string{"-r", "30"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			got1, got2, err := tc.in.Args()

			if tc.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected error: %s but no error was returned", tc.expectedErr)
				}
				if err.Error() != tc.expectedErr {
					t.Fatalf("error messages do not match, got: %s, expect: %s", err.Error(), tc.expectedErr)
				}
			} else if err != nil && tc.expectedErr == "" {
				t.Fatalf("got unexpected error: %s", err.Error())
			}
			if tc.expectedInputArgs != nil {
				if diff := cmp.Diff(got1, tc.expectedInputArgs); diff != "" {
					t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
				}
			} else {
				if diff := cmp.Diff(got2, tc.expectedArgs); diff != "" {
					t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
				}
			}

		})
	}

}

func TestNewFromInterface(t *testing.T) {
	tcs := []struct {
		name        string
		in          map[interface{}]interface{}
		expected    FfmpegOpts
		expectedErr string
	}{
		{
			name: "default video setting",
			in: map[interface{}]interface{}{
				"name": "bla",
			},
			expected: FfmpegOpts{
				Name:           "bla",
				videoExtension: "mp4",
			},
		},

		{
			name: "Threads",
			in: map[interface{}]interface{}{
				"name":    "bla",
				"threads": "3",
			},
			expected: FfmpegOpts{
				Name:           "bla",
				videoExtension: "mp4",
				Threads:        3,
			},
		},
		{
			name: "codec",
			in: map[interface{}]interface{}{
				"name":  "bla",
				"codec": "libx264",
			},
			expected: FfmpegOpts{
				Name:       "bla",
				VideoCodec: "libx264",
			},
		},
		{
			name: "wrong codec",
			in: map[interface{}]interface{}{
				"name":  "bla",
				"codec": "bla",
			},
			expectedErr: NotAllowedCodec,
		},
		{
			name: "quality",
			in: map[interface{}]interface{}{
				"name":           "bla",
				"quality_crf":    "23",
				"quality_preset": "medium",
				"quality_tune":   "film",
			},
			expected: FfmpegOpts{
				Name:          "bla",
				QualityCRF:    getIntPointer(23),
				QualityPreset: "medium",
				QualityTune:   "film",
			},
		},
		{
			name: "scale",
			in: map[interface{}]interface{}{
				"name":  "bla",
				"scale": 480,
			},
			expected: FfmpegOpts{
				Name:       "bla",
				VideoScale: 480,
			},
		},

		{
			name: "scale string",
			in: map[interface{}]interface{}{
				"name":  "bla",
				"scale": "480",
			},
			expected: FfmpegOpts{
				Name:       "bla",
				VideoScale: 480,
			},
		},

		{
			name: "duration and start",
			in: map[interface{}]interface{}{
				"name":     "bla",
				"duration": "2m",
				"start":    "10s",
			},
			expected: FfmpegOpts{
				Name:          "bla",
				VideoDuration: 2 * time.Minute,
				VideoStart:    10 * time.Second,
			},
		},

		{
			name: "extra field",
			in: map[interface{}]interface{}{
				"name":  "bla",
				"extra": "-r 30",
			},
			expected: FfmpegOpts{
				Name:  "bla",
				Extra: "-r 30",
			},
		},
		{
			name: "cuda decoding",
			in: map[interface{}]interface{}{
				"name":           "cuda decoding",
				"cuda_decoding":  "true",
				"cuda_hw_output": "true",
			},
			expected: FfmpegOpts{
				Name:         "cuda decoding",
				CudaDecoding: true,
				CudaHwOutput: true,
			},
		},
		{
			name: "cuda decoding disabled",
			in: map[interface{}]interface{}{
				"name":           "cuda decoding",
				"cuda_decoding":  "false",
				"cuda_hw_output": "",
			},
			expected: FfmpegOpts{
				Name:         "cuda decoding",
				CudaDecoding: false,
				CudaHwOutput: false,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			opts, err := NewFromInterface(tc.in)

			if tc.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected error: %s but no error was returned", tc.expectedErr)
				}
				if err.Error() != tc.expectedErr {
					t.Fatalf("error messages do not match, got: %s, expect: %s", err.Error(), tc.expectedErr)
				}
			} else if err != nil && tc.expectedErr == "" {
				t.Fatalf("got unexpected error: %s", err.Error())
			} else {
				if diff := cmp.Diff(&tc.expected, opts, cmpopts.IgnoreUnexported(FfmpegOpts{})); diff != "" {
					t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
				}
			}

		})
	}
}
