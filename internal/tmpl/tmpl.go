package tmpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	tmplStr string
}

func NewTmplFromFile(file string) (Template, error) {

	dat, err := os.ReadFile(file)
	if err != nil {
		return Template{}, err
	}

	tmpl := Template{
		tmplStr: string(dat),
	}

	return tmpl, nil
}

type TemplateNotFoundErr struct {
	tmpl string
}

type TemplateData struct {
	Args    []string `json:"args"`
	FileExt string   `json:"extension"`
}

func (tmpl Template) Parse(data any) (TemplateData, error) {
	tpl := tmpl.tmplStr
	tpl = strings.ReplaceAll(tpl, "\n", " ")

	t, err := template.New("irrelevant").Parse(tpl)
	if err != nil {
		return TemplateData{}, fmt.Errorf("unable to parse template: %s", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return TemplateData{}, err
	}

	td := TemplateData{}
	err = json.Unmarshal(buf.Bytes(), &td)
	if err != nil {
		return td, fmt.Errorf("unable to unmarshal json template: %s", err)
	}
	td.Args = dropEmpty(td.Args)

	return td, nil
}

// remove empty items in slice
func dropEmpty(in []string) []string {
	var out []string
	for _, v := range in {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}

func (t TemplateNotFoundErr) Error() string {
	return fmt.Sprintf("template \"%s\" not found", t.tmpl)
}

// FindTemplate searches the list of folders for a template file named like the provided name
// returns the path of the template,
// if a template is present in more than one folder, the last one will be returned
func FindTemplate(folders []string, name string) (string, error) {
	fPath := ""
	for _, folder := range folders {
		files, err := os.ReadDir(folder)
		if err != nil {

			if strings.Contains(err.Error(), "no such file or directory") {
				continue
			}
			return "", err
		}

		for _, file := range files {

			if !file.IsDir() && file.Name() == name+".tmpl" {
				fPath = filepath.Join(folder, file.Name())
			}
		}
	}

	if fPath == "" {
		return "", TemplateNotFoundErr{
			tmpl: name,
		}
	}

	return fPath, nil
}