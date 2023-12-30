package tmpl

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFindTemplate(t *testing.T) {

	tcs := []struct {
		name     string
		tmplName string
		folders  []string
		expect   string
		expetErr string
	}{
		{
			name:     "find file",
			tmplName: "tmpl1",
			folders: []string{
				"testdata/templates/folder1",
			},
			expect: "testdata/templates/folder1/tmpl1.tmpl",
		},
		{
			name:     "overlayed file",
			tmplName: "tmpl1",
			folders: []string{
				"testdata/templates/folder1",
				"testdata/templates/folder2",
			},
			expect: "testdata/templates/folder2/tmpl1.tmpl",
		},
		{
			name:     "only existing in last",
			tmplName: "tmpl2",
			folders: []string{
				"testdata/templates/folder1",
				"testdata/templates/folder2",
			},
			expect: "testdata/templates/folder2/tmpl2.tmpl",
		},
		{
			name:     "not found",
			tmplName: "tmpl3",
			folders: []string{
				"testdata/templates/folder1",
				"testdata/templates/folder2",
			},
			expetErr: "template \"tmpl3\" not found",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			got, err := FindTemplate(tc.folders, tc.tmplName)
			if tc.expetErr != "" {
				if err == nil {
					t.Fatal("expecting an error but none returned")
				}
				if err.Error() != tc.expetErr {
					t.Fatalf("expecting error msg: %s, but got instead %s", err.Error(), tc.expetErr)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}

				if diff := cmp.Diff(got, tc.expect); diff != "" {
					t.Errorf("unexpected value (-got +want)\n%s", diff)
				}
			}
		})
	}

}

func TestTemplateArgs(t *testing.T) {

	type data struct {
		Key string
	}

	tcs := []struct {
		name     string
		tmplFile string
		data     data
		expect   TemplateData
	}{
		{
			name:     "find file",
			tmplFile: "testdata/templates/tc_simple.tmpl",
			data: data{
				Key: "SomeValue",
			},
			expect: TemplateData{
				Args: []string{
					"-v",
					"-key",
					"value",
				},
				FileExt: "mkv",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tmpl, err := NewTmplFromFile(tc.tmplFile)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got, err := tmpl.Parse(tc.data)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if diff := cmp.Diff(got, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}

}
