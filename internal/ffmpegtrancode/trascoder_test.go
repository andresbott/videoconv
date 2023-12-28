package ffmpegtrancode

import (
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"testing"
)

func TestGetCmd(t *testing.T) {

	inputFile := "testdata/video.mp4"
	outputFile := "output.mp4"

	tcs := []struct {
		name   string
		in     Cfg
		tmpl   Template
		expect []string
	}{
		{
			name: "simple",
			in: Cfg{
				FfmpegBin: "/usr/bin/ffmpeg",
			},
			tmpl: Template{
				tmplStr: "-v -key value",
			},
			expect: []string{
				"/usr/bin/ffmpeg",
				"-i",
				filepath.Join(absPath(), inputFile),
				"-v",
				"-key",
				"value",
				filepath.Join(absPath(), outputFile),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ffmpeg, err := New(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			cmd, err := ffmpeg.GetCmd(inputFile, outputFile, tc.tmpl)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(cmd, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}

func absPath() string {
	abs, _ := filepath.Abs("./")
	return abs
}
