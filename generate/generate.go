package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// TemplateMap is used to store templates
var TemplateMap map[string]string //Ussed to store templates

// Generate is a type
type Generate struct {
	Root       string // ideally project root directory .i.e project name
	DbType     string
	HasHandler bool
	Models     []Model
	Templates  map[string]string
}

// Model is to create a model
type Model struct {
	Name   string
	Fields []Field
}

// Field is to create a field
type Field struct {
	Name string
	Type string
}

// New is to generate a new template
func New(file string) (tg *Generate, err error) {
	ext := filepath.Ext(file)
	fmt.Println(ext)
	if ext != ".json" {
		return nil, errors.New("Only json files are allowed ")
	}
	cFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cFile), &tg)
	if err != nil {
		return nil, err
	}
	err = tg.MkDirs() // Generate required directories
	if err != nil {
		return nil, err
	}

	return tg, nil
}

// LoadTemplates is to load templates from template folder to the map
func LoadTemplates(path string) (err error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	for _, file := range files {

		content, err := ioutil.ReadFile("templates/" + file.Name())

		if err != nil {
			return err
		}

		TemplateMap[file.Name()] = string(content)
	}
	return nil
}

// GenerateAllModelFiles is to create all model files
func (tg *Generate) GenerateAllModelFiles(tmpl string) (err error) {
	for _, v := range tg.Models {
		modelsFile := path.Join(tg.Root, "models", v.Name+".go")
		err = TmplToFile(modelsFile, tmpl, v)
		if err != nil {
			return err
		}

	}
	return err
}

// MkDirs Create all required directories
func (tg *Generate) MkDirs() (err error) {
	if tg == nil || tg.Root == "" {
		return errors.New("project root directory or the generation has error")
	}
	err = os.Mkdir(tg.Root, 0777)
	if err != nil {
		return err
	}

	models := filepath.Join(tg.Root, "models")
	err = os.Mkdir(models, 0777)
	if err != nil {
		return err
	}

	if tg.HasHandler {
		handlers := filepath.Join(tg.Root, "handlers")
		err = os.Mkdir(handlers, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateMain main.go
func (tg *Generate) CreateMain(tmpl string) (err error) {
	fileName := path.Join(tg.Root, "main.go")

	data := make(map[string]string)

	data["project_name"] = tg.Root

	err = TmplToFile(fileName, tmpl, data)
	if err != nil {
		return err
	}
	return nil
}

// TmplToString is to convert from tmpl to a string
func TmplToString(tmpl string, data interface{}) (result string, err error) {
	t := template.Must(template.New("toString").Parse(tmpl))
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)

	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

// TmplToFile is to convert from template to a file
func TmplToFile(filePath string, tmpl string, data interface{}) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	t := template.Must(template.New("toFile").Parse(tmpl))

	err = t.Execute(file, data)

	if err != nil {
		return err
	}

	return nil
}

// Generater interface is to provide generater methods
type Generater interface {
	TmplToString(tmpl string, data interface{}) (result string, err error)
	TmplToFile(filePath string, tmpl string, data interface{}) (err error)
}
