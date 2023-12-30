package videoconv

import (
	"github.com/AndresBott/videoconv/app/videoconv/config"
	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func newVideConv(t *testing.T) (*Converter, string) {

	// prepare the stage
	tmpDir := t.TempDir()
	location := filepath.Join(tmpDir, "sample")
	dirs := []string{
		"sample/in/nested",
		"sample/out",
		"sample/tmp",
		"sample/templates",
		"sample/fail",
	}
	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, d), 0755)
		if err != nil {
			t.Fatalf("error preparing the stage: %v", err)
		}
	}

	input, err := os.ReadFile("testdata/video.mp4")
	if err != nil {
		t.Fatalf("error preparing the stage: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "sample/in/nested/video.mp4"), input, 0644)
	if err != nil {
		t.Fatalf("error preparing the stage: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "sample/in/video1.MKV"), []byte("content"), 0644)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// templates
	templates := map[string]string{
		"empty":        `{"args":[]}`,
		"mkv":          `{"args":[],"extension":"mkv"}`,
		"broken-param": `{"args":["-a"]}`,
	}

	for k, v := range templates {
		err = os.WriteFile(filepath.Join(tmpDir, "sample/templates/"+k+".tmpl"), []byte(v), 0644)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	cfg := config.Conf{
		LogLevel: "error",
		TmplDirs: []string{
			filepath.Join(tmpDir, "sample/templates"),
		},

		FfmpegPath:  config.DefaultFFmpeg,
		FfprobePath: config.DefaultFFprobe,
		VideoExtensions: []string{
			"avi", "mkv", "mov", "mp4",
		},
		ConfigLocation: filepath.Join(tmpDir, "videoconv.yaml"),
		Locations: []config.Location{
			{
				Path:      "sample",
				InputDir:  "in",
				OutputDir: "out",
				TmpDir:    "tmp",
				FailDir:   "fail",
			},
		},
	}
	vc, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return vc, location
}

func TestCheck(t *testing.T) {

	vc, locationDir := newVideConv(t)

	// check but don`t create
	err := vc.Check(false)
	if err != nil {
		expectedErrMsg := "directory \"" + locationDir + "/out\" does not exits"
		if err.Error() != expectedErrMsg {
			t.Errorf("unexpected error: %v", err)
		}
	}

	// check but this time also create
	err = vc.Check(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(locationDir)
	if err != nil {
		t.Fatal(err)
	}

	got := []string{}
	for _, e := range entries {
		got = append(got, e.Name())
	}

	want := []string{
		"fail",
		"in",
		"out",
		"templates",
		"tmp",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}
}

func TestRunLocation(t *testing.T) {
	vc, tmpPath := newVideConv(t)
	err := vc.Check(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	type result struct {
		Videos []string
		In     string
		Out    string
		Tmp    string
		Fail   string
	}
	want := result{
		Videos: []string{
			filepath.Join(tmpPath, "in", "nested/video.mp4"),
			filepath.Join(tmpPath, "in", "video1.MKV"),
		},
		In:   filepath.Join(tmpPath, "in"),
		Out:  filepath.Join(tmpPath, "out"),
		Tmp:  filepath.Join(tmpPath, "tmp"),
		Fail: filepath.Join(tmpPath, "fail"),
	}

	got := result{}
	processFn = func(absVideo, absIn, absOut, absTmp, absFail string, profiles []config.Profile) {
		got.Videos = append(got.Videos, absVideo)
		got.In = absIn
		got.Out = absOut
		got.Tmp = absTmp
		got.Fail = absFail
	}

	vc.runLocation(vc.Cfg.Locations[0])

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}
}

func TestProcessVideo(t *testing.T) {

	tcs := []struct {
		name     string
		profiles []config.Profile
		prepare  func(path string, t *testing.T)
		expect   []string
	}{
		{
			name: "simple conversion",
			profiles: []config.Profile{
				{
					Name:     "test",
					Template: "empty",
					Args: map[string]string{
						"key": "value",
					},
				},
			},
			expect: []string{
				"in/video1.MKV",
				"out/nested/video.mp4",
				"out/nested/video.test.mp4",
			},
		},
		{
			name: "delete old tmp before conversion",
			profiles: []config.Profile{
				{
					Name:     "test",
					Template: "empty",
					Args: map[string]string{
						"key": "value",
					},
				},
			},
			expect: []string{
				"in/video1.MKV",
				"out/nested/video.mp4",
				"out/nested/video.test.mp4",
			},
			prepare: func(path string, t *testing.T) {
				// generate a tmp file simulating a leftover
				err := os.WriteFile(filepath.Join(path, "tmp/video.test.mp4"), []byte("content"), 0644)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},

		{
			name: "ffmpeg failure",
			profiles: []config.Profile{
				{
					Name:     "test",
					Template: "broken-param",
					Args: map[string]string{
						"key": "value",
					},
				},
			},
			expect: []string{
				"fail/nested/video.mp4",
				"in/video1.MKV",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			log.SetOutput(io.Discard)

			vc, tmpPath := newVideConv(t)
			videoPath := filepath.Join(tmpPath, "in/nested/video.mp4")
			location := vc.Cfg.Locations[0]
			location.Profiles = tc.profiles

			if tc.prepare != nil {
				tc.prepare(tmpPath, t)
			}

			inPath := filepath.Join(tmpPath, "in")
			outPath := filepath.Join(tmpPath, "out")
			tmp := filepath.Join(tmpPath, "tmp")
			fail := filepath.Join(tmpPath, "fail")

			vc.processVideo(videoPath, inPath, outPath, tmp, fail, location.Profiles)

			files := []string{}
			err := filepath.Walk(tmpPath, func(fPath string, fInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if fInfo.IsDir() {
					return nil
				}

				rel, err := filepath.Rel(tmpPath, fPath)
				if err != nil {
					return err
				}
				ext := filepath.Ext(fPath)
				ext = ext[1:]
				videoEx := []string{
					"mp4", "mkv",
				}
				if !isVideo(ext, videoEx) {
					return nil
				}

				files = append(files, rel)
				return nil
			})
			if err != nil {
				t.Fatalf("unexpected error while walking test directory:%v", err)
			}

			if diff := cmp.Diff(files, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}

}
