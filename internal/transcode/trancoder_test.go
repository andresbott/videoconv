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

	smplOpts := FfmpegOpts{
		VideoCodec: "libx264",
	}

	tcs := []struct {
		name     string
		bin      string
		in       string
		out      string
		opts     FfmpegOpts
		expected string
	}{
		{
			name:     "happyPathMainConf",
			bin:      "/bin/ffmpeg",
			in:       "/abs/path/smpl charsß.mp4",
			opts:     smplOpts,
			out:      "/abs/path/out/smpl charsß.mp4",
			expected: `/bin/ffmpeg -i "/abs/path/smpl charsß.mp4" -c:v libx264 "/abs/path/out/smpl charsß.mp4"`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			cfg := Cfg{
				FfmpegBin:  tc.bin,
				FfmpegOpts: tc.opts,
				InputFile:  tc.in,
				OutputFile: tc.out,
			}
			transcoder, err := New(&cfg)
			if err != nil {
				t.Fatal(err)
			}

			cmd, err := transcoder.GetCmd()
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
