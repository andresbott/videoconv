package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var _ = spew.Dump

func TestLoadConfig(t *testing.T) {

	tcs := []struct {
		name   string
		in     string
		expect Conf
	}{
		{
			name: "defaults",
			in:   "testdata/defaults.yaml",
			expect: Conf{
				LogLevel:        "info",
				Sleep:           5 * time.Minute,
				FfmpegPath:      "/usr/bin/ffmpeg",
				FfprobePath:     "/usr/bin/ffprobe",
				VideoExtensions: []string{"avi", "mkv", "mov"},
				Locations: []Location{
					{
						Path:      "./",
						InputDir:  "in",
						OutputDir: "out",
						TmpDir:    "tmp",
						FailDir:   "fail",
						Profiles:  nil,
					},
				},
				TmplDirs: []string{
					"/etc/videconv/templates",
					"./sample/templates",
				},
			},
		},
		{
			name: "all settings",
			in:   "testdata/allsettings.yaml",
			expect: Conf{
				LogLevel:        "error",
				Sleep:           10 * time.Second,
				FfmpegPath:      "/usr/local/bin/ffmpeg-static",
				FfprobePath:     "/usr/local/bin/ffprobe-static",
				VideoExtensions: []string{"mkv"},
				Locations: []Location{
					{
						Path:      "./",
						InputDir:  "in",
						OutputDir: "out",
						TmpDir:    "tmp",
						FailDir:   "fail",
						Profiles: []Profile{
							{
								Template: "mp4-x265aac",
								Args: map[string]string{
									"height":  "720",
									"bitrate": "4M",
								},
							},
							{
								Template: "test",
								Args: map[string]string{
									"key": "value",
								},
							},
						},
					},
					{
						Path:      "./some_path",
						InputDir:  "input",
						OutputDir: "output",
						TmpDir:    "temp",
						FailDir:   "error",
						Profiles:  nil,
					},
				},
				TmplDirs: []string{
					"/etc/videconv/templates",
					"./sample/templates",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewFromFile(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tc.expect, cmpopts.IgnoreFields(Conf{}, "ConfigLocation")); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}

func TestGeneratedConfig(t *testing.T) {
	tmpDir := t.TempDir()
	sampleCfg := SampleCfg()
	cfgFile := filepath.Join(tmpDir, "videconv.yaml")
	err := os.WriteFile(cfgFile, []byte(sampleCfg), 0644)
	if err != nil {
		t.Fatal()
	}

	got, err := NewFromFile(cfgFile)
	if err != nil {
		t.Fatal(err)
	}

	expect := Conf{
		LogLevel:        "info",
		Sleep:           5 * time.Minute,
		FfmpegPath:      "/usr/bin/ffmpeg",
		FfprobePath:     "/usr/bin/ffprobe",
		VideoExtensions: []string{"avi", "mkv", "mov", "wmv", "mp4"},
		ConfigLocation:  cfgFile,
		Locations: []Location{
			{
				Path:      "./sample",
				InputDir:  "in",
				OutputDir: "out",
				TmpDir:    "tmp",
				FailDir:   "fail",
				Profiles: []Profile{
					{
						Template: "sample",
						Args: map[string]string{
							"key": "value",
						},
					},
				},
			},
		},
		TmplDirs: []string{
			"/etc/videconv/templates",
			"./sample/templates",
		},
	}

	if diff := cmp.Diff(got, expect); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}
}
