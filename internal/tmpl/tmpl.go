package tmpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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

// https://stackoverflow.com/questions/44675087/golang-template-variable-isset
// use like: {{if (isset "Email" .)}}
func isset(name string, data interface{}) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}
	return v.FieldByName(name).IsValid()
}

func (tmpl Template) ParseJson(data, target any) error {
	tpl := tmpl.tmplStr
	tpl = strings.ReplaceAll(tpl, "\n", " ")

	funcMap := sprig.FuncMap()
	funcMap["isset"] = isset

	t, err := template.New("videoTmpl").Funcs(funcMap).Parse(tpl)
	if err != nil {
		return fmt.Errorf("unable to parse template: %s", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}
	err = json.Unmarshal(buf.Bytes(), &target)
	if err != nil {
		return fmt.Errorf("unable to unmarshal json template: %s", err)
	}

	return nil
}

const tmplExt = ".tmpl.json"

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

			if !file.IsDir() && file.Name() == name+tmplExt {
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
