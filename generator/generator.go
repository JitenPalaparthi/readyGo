package generator

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

// Generator is a type
type Generator struct {
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
func New(file string) (tg *Generator, err error) {
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

func CreateModel(data Model) (string, error) {
	tmplModel := `package models
	
	type {{ .Name }} struct {
	{{ range .Fields }}
	{{ .Name }} {{ .Type }}
	{{end}}
	}`

	t := template.Must(template.New("models").Parse(tmplModel))
	buf := bytes.NewBufferString("")
	err := t.Execute(buf, data)

	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (tg *Generator) CreateModelFile(file, data string) error {
	modelsFile := path.Join(tg.Root, "models", file+".go")
	f, err := os.Create(modelsFile)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}

func (tg *Generator) GenerateAllModelFiles() (err error) {
	for _, v := range tg.Models {

		data, err := CreateModel(v)
		if err != nil {
			return err
		}
		err = tg.CreateModelFile(v.Name, data)
		if err != nil {
			return err
		}

	}
	return err
}

// MkDirs Create all required directories
func (tg *Generator) MkDirs() (err error) {
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
func (tg *Generator) CreateMain(tmpl string) (err error) {
	fileName := path.Join(tg.Root, "main.go")

	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	t := template.Must(template.New("main").Parse(tmpl))

	data := make(map[string]string)

	data["project_name"] = tg.Root

	err = t.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
