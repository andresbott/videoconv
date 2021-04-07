package videconv

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		expected App
	}{
		{
			name: "happyPath",
			file: "testdata/happyPath.yaml",
			expected: App{
				ConfigFile:   "testdata/happyPath.yaml",
				OverlayFname: "videoconv.yaml",
				locations: []location{
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
				logLevel:   "error",
				threads:    2,
				sleep:      10 * time.Second,
				ffmpegBin:  "/usr/bin/ffmpeg",
				ffProbeBin: "/usr/bin/ffprobe",
				videoExtensions: []string{
					"mp4", "wmv", "mkv",
				},
				profiles: map[string]profile{
					"item2": {
						template: "item2 template",
						name:     "item2",
					},
					"minimalist": {
						template:  "minimalist template",
						name:      "minimalist",
						extension: "mp4",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			app := App{
				ConfigFile: tc.file,
			}
			err := app.loadConfig()

			if err != nil && err.Error() != "video settings not defined" {
				t.Fatal(err)
			}

			if diff := cmp.Diff(app, tc.expected, cmp.AllowUnexported(App{}, location{}), cmpopts.IgnoreUnexported(profile{})); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}
		})
	}
}
