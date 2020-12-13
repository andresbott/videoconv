package config2

import (
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"testing"
	"time"
)

func getCurrentAbsPath() string {
	abs, _ := filepath.Abs("./")
	return abs
}

func TestConfHandler_Load(t *testing.T) {

	tcs := []struct {
		name     string
		file     string
		expected ConfHandler
	}{
		{
			name: "happyPathMainConf",
			file: "testdata/main.yaml",
			expected: ConfHandler{
				locations: []Location{
					{
						path:      getCurrentAbsPath(),
						inputDir:  "input",
						outputDir: "output",
						tmpDir:    "tmpdir",
						failDir:   "faildir",
					},
					{
						path:            getCurrentAbsPath() + "/location1",
						inputDir:        "in",
						outputDir:       "out",
						tmpDir:          "tmp",
						failDir:         "fail",
						appliedProfiles: []string{"item2"},
					},
				},
				logLevel:     "info",
				threads:      2,
				pollInterval: 10 * time.Second,
				ffmpegBin:    "/usr/bin/ffmpeg",
				profiles: map[string]Profile{
					"minimalist": {
						name:      "minimalist",
						extension: "mp4",
						codec:     "h264",
					},
					"item2": {
						name:      "item2",
						extension: "mp4",
						codec:     "h264",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			cfg := ConfHandler{}
			err := cfg.Load(tc.file)

			if err != nil && err.Error() != "video settings not defined" {
				t.Fatal(err)
			}

			if diff := cmp.Diff(cfg, tc.expected, cmp.AllowUnexported(ConfHandler{}, Location{}, Profile{})); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}

		})
	}

}
