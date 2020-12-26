package videconv

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestLocation(t *testing.T) {

	tcs := []struct {
		name         string
		in           map[interface{}]interface{}
		overlayFname string
		expected     *location
	}{
		{
			name: "default location",
			in: map[interface{}]interface{}{
				"base_path": "./",
				"applied":   []interface{}{"a", "b"},
			},
			expected: &location{
				path:            getCurrentAbsPath(),
				inputDir:        "in",
				outputDir:       "out",
				tmpDir:          "tmp",
				failDir:         "fail",
				appliedProfiles: []string{"a", "b"},
			},
		},
		{
			name:         "overlay location",
			overlayFname: "overlay.yaml",
			in: map[interface{}]interface{}{
				"base_path": "testdata/location1",
				"applied":   []interface{}{"a", "b"},
			},
			expected: &location{
				path:            getCurrentAbsPath() + "/testdata/location1",
				inputDir:        "input",
				outputDir:       "output",
				tmpDir:          "tmpdir",
				failDir:         "faildir",
				appliedProfiles: []string{"a", "b", "720_h265_sample"},
			},
		},
		{
			name:         "overlay drop profiles",
			overlayFname: "overlay_drop_applied.yaml",
			in: map[interface{}]interface{}{
				"base_path": "testdata/location1",
				"applied":   []interface{}{"a", "b"},
			},
			expected: &location{
				path:            getCurrentAbsPath() + "/testdata/location1",
				inputDir:        "input",
				outputDir:       "output",
				tmpDir:          "tmpdir",
				failDir:         "faildir",
				appliedProfiles: []string{"720_h265_sample"},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			got, err := newLocation(tc.in)
			if err != nil {
				t.Fatal(err)
			}

			got, err = got.loadOverlay(tc.overlayFname)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tc.expected, cmp.AllowUnexported(location{})); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}

		})
	}
}

func TestOverlayIsCopy(t *testing.T) {

	// create a location instance
	a, err := newLocation(map[interface{}]interface{}{
		"base_path": "./",
		"applied":   []interface{}{"a", "b"},
	})
	if err != nil {
		t.Fatal(err)
	}
	a.appliedProfiles = []string{"a"}

	// call overlay
	b, _ := a.loadOverlay("")

	// modify value of copy
	b.appliedProfiles = []string{"b"}

	// check that original has not changed
	if a.appliedProfiles[0] != "a" {
		t.Fatal("location.Overlay() does not return a copy of the location struct")
	}
}
