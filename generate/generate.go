package generate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"readyGo/generate/configure"
	"strings"

	"golang.org/x/lint"
)

// Generater interface is to provide generater methods
type Generater interface {
	ToString(tmpl string, data interface{}) (result string, err error)
	ToFile(filePath string, tmpl string, data interface{}) (err error)
	Read(key string) interface{}
}

// Configurator interface is to fetch configration related things
type Configurator interface {
	ReadFC(key string) []configure.StaticFileSet
	ReadTC(key string) []string
	ReadDC(key string) []string
}

// Generate is a type
type Generate struct {
	Project string // ideally project root directory .i.e project name
	Type    string // Type of the project http , grpc , CloudEvents , cli
	Port    string // Port that is used to communicate http project
	DBType  string // mongo , sql based postgres mariadb etc
	Models  []Model
	Gen     Generater
	Con     Configurator
}

// Model is to create a model
type Model struct {
	Name   string
	Fields []Field
}

// Field is to create a field
type Field struct {
	Name        string
	Type        string
	IsKey       bool
	ValidateExp string // Regular expression
}

// New is to generate a new template
func New(file *string, gen Generater, con Configurator) (tg *Generate, err error) {
	if file == nil || *file == "" {
		return nil, errors.New("no file provided")
	}
	if gen == nil {
		return nil, errors.New("template cannot be nil.Load template before creating Generater")
	}
	if con == nil {
		return nil, errors.New("template Confirator cannot be nil.Load Configue before creating Configurator")
	}
	ext := filepath.Ext(*file)
	fmt.Println(ext)
	if ext != ".json" {
		return nil, errors.New("Only json files are allowed ")
	}
	cFile, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cFile), &tg)
	if err != nil {
		return nil, err
	}
	err = tg.ValidateAndChangeIdentifier()
	if err != nil {
		return nil, err
	}

	root := strings.ToLower(tg.Project)

	tg.Project = root

	tg.Gen = gen // Assign generater template loading engine interface

	tg.Con = con // Assign generater configation loading enginer interface

	err = tg.Validate()

	if err != nil {
		return nil, err
	}

	return tg, nil
}

// RmDir is to remove dirs
func (tg *Generate) RmDir() (err error) {
	err = os.RemoveAll(tg.Project)
	if err != nil {
		return err
	}
	return nil
}

func (tg *Generate) ValidateAndChangeIdentifier() (err error) {
	for i, m := range tg.Models {
		tmpModel := m.Name
		tg.Models[i].Name, err = checkName(m.Name)
		if err != nil {
			return errors.New("There is an error at the following model config name in json file:" + tmpModel)
			fmt.Println("----1", err)
			return err
		}
		tg.Models[i].Name = strings.ToUpper(string(tmpModel[0])) + string(tmpModel[1:])

		for j, f := range m.Fields {
			tmpField := f.Name

			tg.Models[i].Fields[j].Name, err = checkName(f.Name)
			if err != nil {
				return errors.New("There is an error at the following field name in json file:" + tmpField)
			}
			tg.Models[i].Fields[j].Name = strings.ToUpper(string(tmpField[0])) + string(tmpField[1:])
		}
	}
	return nil
}

func (tg *Generate) GenerateAll(key string) (err error) {
	if tg == nil {
		return errors.New("temp	late generater is not a valid object.Try to instantiate it through generater.New function")
	}
	if tg.Project == "" {
		return errors.New("invalid root directory")
	}
	// Todo write more conditions here

	// Step-1 create all directories

	dirs := tg.Con.ReadDC(key)

	for _, dir := range dirs {
		path := filepath.Join(tg.Project, dir)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Step-2 copy all static files

	fcs := tg.Con.ReadFC(key)

	for _, fc := range fcs {
		dst := filepath.Join(tg.Project, fc.Dst)
		err = CopyStaticFile(fc.Src, dst)
		if err != nil {
			return err
		}
	}

	// Step-3 generate all files based on templates

	tcs := tg.Con.ReadTC(key)
	mhandler := make(map[string]interface{})
	mhandler["Project"] = tg.Project
	//mhandler["dbType"] = tg.DBType
	mhandler["config"] = tg
	for _, v := range tg.Models {
		mhandler["Model"] = v
		for _, tc := range tcs {
			paths := path.Join(tg.Project, tc, strings.ToLower(v.Name)+".go")
			err = GenerateFile(tg.Gen, mhandler, tc, paths)
			if err != nil {
				return err
			}
		}
	}

	err = tg.CreateMain(tg.Gen.Read("main").(string)) // Generate main.go
	if err != nil {
		return err
	}
	return nil

}

// CreateMain main.go
func (tg *Generate) CreateMain(tmpl string) (err error) {
	fileName := path.Join(tg.Project, "main.go")

	data := make(map[string]string)

	data["project_name"] = "demo" //tg.Models

	if tg.Gen != nil {
		err = tg.Gen.ToFile(fileName, tmpl, *tg)
		if err != nil {
			return err
		}
	} else {
		errors.New("Gen interface is not assigned to Generate object")
	}
	return nil
}

// Validate is to validate the object
func (tg *Generate) Validate() (err error) {
	if tg.Project == "" {
		return errors.New("Project name is missing")
	}

	if tg.Type == "" || (tg.Type != "http" && tg.Type != "grpc" && tg.Type != "cloudEvent" && tg.Type != "cli") {
		return errors.New(" Project type must be http | grpc | cloudEvent | cli")
	}
	if tg.DBType == "" || (tg.DBType != "mongo" && tg.DBType != "sql") {
		return errors.New(" Databas type (DBType) must be mongo | sql ")
	}
	// checking duplicate models and fields
	modelMap := make(map[string]string)
	for _, m := range tg.Models {
		_, ok := modelMap[strings.ToLower(m.Name)]
		if ok {
			return errors.New(" Duplicate model names:" + m.Name)
		}
		modelMap[strings.ToLower(m.Name)] = "noted"
		fieldMap := make(map[string]string)
		for _, f := range m.Fields {
			_, ok := fieldMap[strings.ToLower(f.Name)]
			if ok {
				return errors.New(" Duplicate field names:" + f.Name)
			}
			fieldMap[strings.ToLower(f.Name)] = "noted"
		}
	}

	return err
}

// GenerateFile is to create all model files
func GenerateFile(gen Generater, model map[string]interface{}, key, filePath string) (err error) {
	if key == "" {
		return errors.New("empty template key provided")
	}
	tmpl := gen.Read(key)
	err = gen.ToFile(filePath, tmpl.(string), model)
	if err != nil {
		return err
	}

	return nil
}

// CopyStaticFile is to copy static files
func CopyStaticFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	//err = copyFileContents(src, dst)
	return
}

func checkName(s string) (string, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Package main is awesome\npackage main\n// %s is wonderful\nvar %s int\n", s, s)
	var l lint.Linter
	problems, err := l.Lint("", buf.Bytes())
	if err != nil {
		return "", err
	}
	if len(problems) >= 1 {
		t := problems[0].Text
		if i := strings.Index(t, " should be "); i >= 0 {
			return t[i+len(" should be "):], nil
		}
	}
	return "", nil
}
