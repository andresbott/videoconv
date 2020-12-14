package transcoder

import (
	"github.com/google/go-cmp/cmp"
	"path/filepath"
	"strings"
	"testing"
)

func getCurrentAbsPath() string {
	abs, _ := filepath.Abs("./")
	return abs
}

func TestConfHandler_Load(t *testing.T) {

	tcs := []struct {
		name     string
		bin      string
		in       string
		out      string
		tmp      string
		expected string
	}{
		{
			name:     "happyPathMainConf",
			bin:      "/bin/ffmpeg",
			in:       "/abs/path/sample.mp4",
			out:      "/abs/path/out/",
			tmp:      "/abs/path/out",
			expected: "bla",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			cfg := Cfg{
				FfmpedBin: tc.bin,
				VideoFile: tc.in,
				OutputDir: tc.out,
				TmpDir:    tc.tmp,
			}
			transcoder, err := New(&cfg)
			if err != nil {
				t.Fatal(err)
			}

			cmd, err := transcoder.getCmd()
			if err != nil {
				t.Fatal(err)
			}

			got := strings.Join(cmd, " ")

			if diff := cmp.Diff(got, tc.expected, cmp.AllowUnexported(Transcoder{})); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}

		})
	}

}
