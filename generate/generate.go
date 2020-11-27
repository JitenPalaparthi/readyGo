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
	"sort"
	"strings"

	"golang.org/x/lint"
	"gopkg.in/yaml.v2"
)

// Generater interface is to provide generater methods
type Generater interface {
	ToString(tmpl string, data interface{}) (result string, err error)
	ToFile(filePath string, tmpl string, data interface{}) (err error)
	Read(key string) interface{}
}

// Configurator interface is to fetch configration related things
type Configurator interface {
	ReadFC(key string) []configure.CopyLoc
	ReadTC(key string) []configure.CopyLoc
	ReadSF(key string) []configure.CopyLoc
	ReadDC(key string) []string
}

// Generate is a type
type Generate struct {
	Version string  `json:"version" yaml:"version"`
	Project string  `json:"project" yaml:"project"` // ideally project root directory .i.e project name
	Type    string  `json:"type" yaml:"type"`       // Type of the project http , grpc , CloudEvents , cli
	Port    string  `json:"port" yaml:"port"`       // Port that is used to communicate http project
	DB      string  `json:"db" yaml:"db"`           // mongo , sql based postgres mariadb etc
	Models  []Model `json:"models" yaml:"models"`
	Gen     Generater
	Con     Configurator
}

// Model is to create a model
type Model struct {
	Name   string  `json:"name" yaml:"name"`
	Fields []Field `json:"fields" yaml:"fields"`
}

// Field is to create a field
type Field struct {
	Name        string `json:"name" yaml:"name"`               // Name of the field. Should be valid Go identifier
	Type        string `json:"type" yaml:"type"`               // Go basic types are only allowed
	IsKey       bool   `json:"isKey" yaml:"isKey"`             // If it is a key field . Key fields generates different methods to check the data in the database is unique or not
	ValidateExp string `json:"validateExp" yaml:"validateExp"` // Regular expression that would be used for field level validations in the models
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
	if ext != ".json" && ext != ".yaml" && ext != ".yml" {
		return nil, errors.New("Only json | yaml | yml files are allowed ")
	}
	cFile, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, err
	}
	if ext == ".json" {
		err = json.Unmarshal([]byte(cFile), &tg)
		if err != nil {
			return nil, err
		}
	}

	if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal([]byte(cFile), &tg)
		if err != nil {
			return nil, err
		}
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
	mhandler["config"] = tg
	for _, v := range tg.Models {
		mhandler["Model"] = v
		for _, tc := range tcs {
			paths := path.Join(tg.Project, tc.Dst, strings.ToLower(v.Name)+".go")
			err = GenerateFile(tg.Gen, mhandler, tc.Src, paths)
			if err != nil {
				return err
			}
		}
	}

	// Step-4 generate one and only files based on templates
	sf := tg.Con.ReadSF(key)
	mhandler["config"] = tg
	for _, f := range sf {
		paths := path.Join(tg.Project, f.Dst)
		err = GenerateFile(tg.Gen, mhandler, f.Src, paths)
		if err != nil {
			return err
		}
	}

	/*err = tg.CreateMain(tg.Gen.Read("main").(string)) // Generate main.go
	if err != nil {
		return err
	}*/
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
	if tg.DB == "" || (tg.DB != "mongo" && tg.DB != "sql") {
		return errors.New(" Databas type (DB) must be mongo | sql ")
	}
	// The following are supported types
	basicSupportedTypes := []string{"bool", "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64"}

	// checking duplicate models and fields
	modelMap := make(map[string]string)
	for i, m := range tg.Models {
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

			// Validate field types
			if !strings.Contains(strings.Join(basicSupportedTypes, ","), strings.ToLower(f.Type)) {
				return errors.New("Not supported field type:" + f.Type)
			}
		}
		if tg.DB == "mongo" {
			_, ok := fieldMap["id"]
			if !ok {
				id := Field{Name: "Id", Type: "string"}
				tg.Models[i].Fields = append(tg.Models[i].Fields, id)
			}
		}
	}

	return err
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
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
