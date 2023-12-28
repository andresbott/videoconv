package videoconv

import (
	"github.com/AndresBott/videoconv/app/videoconv/config"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func newVideConv(t *testing.T) (*Converter, string) {
	tmpDir := t.TempDir()
	location := filepath.Join(tmpDir, "banana")
	err := os.Mkdir(location, 0755)
	if err != nil {
		t.Errorf("unable to create dir %s, %v", location, err)
	}

	// populate with fake data
	inputDir := filepath.Join(location, "in")
	err = os.Mkdir(inputDir, 0755)
	if err != nil {
		t.Errorf("unable to create dir %s, %v", inputDir, err)
	}

	content := []byte("content")
	err = os.WriteFile(filepath.Join(location, "in", "video1.avi"), content, 0644)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	subDir := filepath.Join(location, "in", "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Errorf("unable to create dir %s, %v", subDir, err)
	}
	nested := filepath.Join(subDir, "nested")
	err = os.Mkdir(nested, 0755)
	if err != nil {
		t.Errorf("unable to create dir %s, %v", nested, err)
	}
	err = os.WriteFile(filepath.Join(nested, "video1.MKV"), content, 0644)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	cfg := config.Conf{
		ConfigLocation: filepath.Join(tmpDir, "videoconv.yaml"),
		LogLevel:       "error",
		Locations: []config.Location{
			{
				Path:      "banana",
				InputDir:  "in",
				OutputDir: "out",
				TmpDir:    "tmp",
				FailDir:   "fail",
			},
		},
		VideoExtensions: []string{
			"avi", "mkv", "mov",
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
		"tmp",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}
}

func TestRunLocation(t *testing.T) {
	vc, _ := newVideConv(t)
	err := vc.Check(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	var got []string
	want := []string{
		"in/subdir/nested/video1.MKV",
		"in/video1.avi",
	}

	vc.runLocation(vc.Cfg.Locations[0], func(video string) {
		got = append(got, video)
	})

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected value (-got +want)\n%s", diff)
	}
}
