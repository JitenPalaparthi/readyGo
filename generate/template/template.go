package template

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// TmplMap is used to store templates
type TmplMap map[string]string

// New is to load templates from template folder to the map
func New(path string) (tm TmplMap, err error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	for _, file := range files {

		content, err := ioutil.ReadFile("templates/" + file.Name())

		if err != nil {
			return nil, err
		}

		if tm == nil {
			tm = make(map[string]string)
		}

		tm[file.Name()] = string(content)
	}
	return tm, nil
}

// ToString is to convert from tmpl to a string
func (tm TmplMap) ToString(tmpl string, data interface{}) (result string, err error) {
	t := template.Must(template.New("toString").Funcs(template.FuncMap{
		"ToLower": func(str string) string {
			return strings.ToLower(str)
		},
	}).Funcs(
		template.FuncMap{
			"Initial": func(str string) string {
				if len(str) > 0 {
					return string(strings.ToLower(str)[0])
				}
				return "x"
			},
		}).Parse(tmpl))
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)

	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

// ToFile is to convert from template to a file
func (tm TmplMap) ToFile(filePath string, tmpl string, data interface{}) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	t := template.Must(template.New("toFile").Funcs(template.FuncMap{
		"ToLower": func(str string) string {
			return strings.ToLower(str)
		},
	}).Funcs(
		template.FuncMap{
			"Initial": func(str string) string {
				if len(str) > 0 {
					return string(strings.ToLower(str)[0])
				}
				return "x"
			},
		}).Parse(tmpl))

	err = t.Execute(file, data)

	if err != nil {
		return err
	}

	return nil
}

// GetTmplMap is to fetch the loaded map
func (tm TmplMap) Read(key string) interface{} {
	return tm[key]
}
