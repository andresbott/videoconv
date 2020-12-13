package config2

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestProfiles(t *testing.T) {

	tcs := []struct {
		name     string
		in       map[interface{}]interface{}
		expected Profile
	}{
		{
			name: "default video setting",
			in: map[interface{}]interface{}{
				"name": "bla",
			},
			expected: Profile{
				name:      "bla",
				extension: "mp4",
				codec:     "h264",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			vst, err := newProfile(tc.in)

			if err != nil && err.Error() != "video settings not defined" {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, vst, cmp.AllowUnexported(ConfHandler{}, Location{}, Profile{})); diff != "" {
				t.Errorf("%s: (-got +want)\n%s", tc.name, diff)
			}

		})
	}

}
