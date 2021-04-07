package videconv

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

// this verifies that a correct error is returned on new profile
func TestNewProfile(t *testing.T) {

	t.Run("type conversion error", func(t *testing.T) {
		_, err := newProfile("someRandom string")

		expect := "unable to cast input from map[interface{}]interface{}"
		if err.Error() != expect {
			t.Fatalf("expected error: \"%s\" but got: \"%s\"", expect, err.Error())
		}
	})

	t.Run("empty template error", func(t *testing.T) {
		payload := map[interface{}]interface{}{
			"name":     "test",
			"template": "",
		}

		_, err := newProfile(payload)
		expect := "template cannot be empty"
		if err.Error() != expect {
			t.Fatalf("expected error: \"%s\" but got: \"%s\"", expect, err.Error())
		}
	})

	t.Run("get default extension", func(t *testing.T) {
		payload := map[interface{}]interface{}{
			"name":     "test",
			"template": "-i",
		}

		expect := profile{
			template:  "-i",
			name:      "test",
			extension: "mp4",
		}

		got, _ := newProfile(payload)

		if diff := cmp.Diff(got, &expect, cmp.AllowUnexported(profile{})); diff != "" {
			t.Errorf("%s: (-got +want)\n%s", "get default extension", diff)
		}
	})
}
