package transcoder

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestConfHandler_Load(t *testing.T) {

	tcs := []struct {
		name string
		bin  string
		in   string
		out  string
		tmpl string

		expected []string
	}{
		{
			name: "happyPathMainConf",
			bin:  "/bin/ffmpeg",
			in:   "/abs/path/smpl charsß.mp4",
			out:  "/abs/path/out/smpl charsß.mp4",
			tmpl: `-i "{{ .Input }}" -o "{{ .Output}}"`,
			expected: []string{
				"/bin/ffmpeg",
				"-i",
				`"/abs/path/smpl charsß.mp4"`,
				"-o",
				`"/abs/path/out/smpl charsß.mp4"`,
			},
		},
		{
			name: "happyPathMainConf",
			bin:  "/bin/ffmpeg",
			in:   "/abs/path/smpl charsß.mp4",
			out:  "/abs/path/out/smpl charsß.mp4",
			tmpl: `-i "{{ .Input }}" 
-o "{{ .Output}}"`,
			expected: []string{
				"/bin/ffmpeg",
				"-i",
				`"/abs/path/smpl charsß.mp4"`,
				"-o",
				`"/abs/path/out/smpl charsß.mp4"`,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			cfg := Cfg{
				FfmpegBin:  tc.bin,
				InputFile:  tc.in,
				OutputFile: tc.out,
				Template:   tc.tmpl,
			}

			transcoder, err := New(&cfg)
			if err != nil {
				t.Fatal(err)
			}

			cmd, err := transcoder.GetCmd()
			if err != nil {
				t.Fatal(err)
			}

			got := cmd

			if diff := cmp.Diff(got, tc.expected, cmp.AllowUnexported(Transcoder{})); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}

		})
	}
}

func TestFileMetadata(t *testing.T) {
	tcs := []struct {
		name     string
		file     string
		expected tmplData
	}{
		{
			name: "probe file 1",
			file: "testdata/kodak_instamatic_320x180.mp4",
			expected: tmplData{
				Input:           "",
				Output:          "",
				VideoWidth:      320,
				VideoHeight:     180,
				DurationSeconds: 2.262,
				FrameRate:       25,
			},
		},
		{
			name: "probe file 2",
			file: "testdata/sample with erroneus name &$·ŋ640x360.webm",
			expected: tmplData{
				Input:           "",
				Output:          "",
				VideoWidth:      640,
				VideoHeight:     360,
				DurationSeconds: 2.243,
				FrameRate:       25,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			transcoder, err := New(&Cfg{
				FfProbeBin: "/usr/bin/ffprobe", // WARNING: unfortunately the test (like the binary) depends on having ffprobe installed
				InputFile:  tc.file,
			})
			if err != nil {
				t.Fatal(err)
			}

			got := tmplData{}
			err = transcoder.scanF(&got)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}
		})
	}
}
