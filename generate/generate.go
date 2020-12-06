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
	"readyGo/mapping"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/lint"
	"gopkg.in/yaml.v2"
)

// Generate is a type that holds configuration data
type Generate struct {
	Version  string   `json:"version" yaml:"version"`
	Project  string   `json:"project" yaml:"project"` // ideally project root directory .i.e project name
	Kind     string   `json:"kind" yaml:"kind"`       // Kind of the project http , grpc , CloudEvents , cli
	Port     string   `json:"port" yaml:"port"`       // Port that is used to communicate http project
	DB       string   `json:"db" yaml:"db"`           // mongo , sql based postgres mariadb etc
	Database Database `json:"database" yaml:"database"`
	Models   []Model  `json:"models" yaml:"models"`
	Mapping  *mapping.Mapping
}

// Model is to hold model data from configuration file
type Model struct {
	Name   string  `json:"name" yaml:"name"`
	Fields []Field `json:"fields" yaml:"fields"`
}

// Field is to hold fields in a model that comes from configuration file
type Field struct {
	Name        string `json:"name" yaml:"name"`               // Name of the field. Should be valid Go identifier
	Type        string `json:"type" yaml:"type"`               // Go basic types are only allowed
	IsKey       bool   `json:"isKey" yaml:"isKey"`             // If it is a key field . Key fields generates different methods to check the data in the database is unique or not
	ValidateExp string `json:"validateExp" yaml:"validateExp"` // Regular expression that would be used for field level validations in the models
}

// Database struct type contains database related information
type Database struct {
	Kind             string `json:"kind" yaml:"kind"`
	ConnectionString string `json:"connectionString" yaml:"connectionString"`
	Name             string `json:"name" yaml:"name"`
}

// New is to generate a new generater.
func New(file *string, mapping *mapping.Mapping) (tg *Generate, err error) {

	if file == nil || *file == "" {
		return nil, errors.New("no file provided")
	}
	if mapping == nil {
		return nil, errors.New("mapping cannot be empty")
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

	tg.Mapping = mapping

	err = tg.Validate()

	if err != nil {
		return nil, err
	}

	return tg, nil
}

// CreateAll creates all kinds of files based on the provided mappings
func (tg *Generate) CreateAll() (err error) {
	if tg == nil {
		return errors.New("temp	late generater is not a valid object.Try to instantiate it through generater.New function")
	}
	if tg.Project == "" {
		return errors.New("invalid root directory")
	}
	// Todo write more conditions here

	for _, opsData := range tg.Mapping.OpsData {
		//fmt.Println(opsData)
		switch opsData.OpType {
		case "directories":
			path := filepath.Join(tg.Project, opsData.Src)
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				errRm := tg.RmDir()
				if errRm != nil {
					return errors.New("1." + err.Error() + ".2." + errRm.Error())
				}
				return err
			}
		case "static-files":
			dst := filepath.Join(tg.Project, opsData.Dst)
			content, err := tg.Mapping.Reader.Read(opsData.Src)
			if err != nil {
				errRm := tg.RmDir()
				if errRm != nil {
					return errors.New("1." + err.Error() + ".2." + errRm.Error())
				}
				return err
			}
			if content != "" {
				var li int

				cos := runtime.GOOS
				switch cos {
				case "windows":
					li = strings.LastIndex(dst, "\\")
				default:
					li = strings.LastIndex(dst, "/")
				}

				dirs := dst[0:li]
				err = os.MkdirAll(dirs, 0755)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}

				err = ioutil.WriteFile(dst, []byte(content), 0644)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}
			}
		case "multiple-file-templates":
			mhandler := make(map[string]interface{})
			mhandler["Project"] = tg.Project
			mhandler["config"] = tg
			for _, v := range tg.Models {
				mhandler["Model"] = v
				dst := filepath.Join(tg.Project, opsData.Dst)
				err = os.MkdirAll(dst, 0755)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}

				dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+".go")
				content, err := tg.Mapping.Reader.Read(opsData.Src)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}
				if content != "" {
					err := WriteTmplToFile(dst, content, mhandler)
					if err != nil {
						errRm := tg.RmDir()
						if errRm != nil {
							return errors.New("1." + err.Error() + ".2." + errRm.Error())
						}
						return err
					}
				}
			}

		case "single-file-templates":
			mhandler := make(map[string]interface{})
			mhandler["config"] = tg
			dst := path.Join(tg.Project, opsData.Dst)
			li := strings.LastIndex(dst, "/")
			dirs := dst[0:li]
			err = os.MkdirAll(dirs, 0755)
			if err != nil {
				errRm := tg.RmDir()
				if errRm != nil {
					return errors.New("1." + err.Error() + ".2." + errRm.Error())
				}
				return err
			}
			content, err := tg.Mapping.Reader.Read(opsData.Src)
			if err != nil {
				errRm := tg.RmDir()
				if errRm != nil {
					return errors.New("1." + err.Error() + ".2." + errRm.Error())
				}
				return err
			}
			if content != "" {
				err := WriteTmplToFile(dst, content, mhandler)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}
			}
		default:
			return errors.New(opsData.OpType + ":this type has no implementation")
		}
	}

	return nil
}

// RmDir is to remove dirs
func (tg *Generate) RmDir() (err error) {
	err = os.RemoveAll(tg.Project)
	if err != nil {
		return err
	}
	return nil
}

// ValidateAndChangeIdentifier is to validate and Change as and where required
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

// Validate is to validate the object
func (tg *Generate) Validate() (err error) {
	if tg.Project == "" {
		return errors.New("Project name is missing")
	}
	if tg.Kind == "" || (tg.Kind != "http" && tg.Kind != "grpc" && tg.Kind != "cloudEvent" && tg.Kind != "cli") {
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
			//	if ok{
			//  todo if the type of the field for id is not string .. it has to be string
			//	}
		}
		if tg.DB == "sql" {
			_, ok := fieldMap["id"]
			if !ok {
				id := Field{Name: "Id", Type: "int"}
				tg.Models[i].Fields = append(tg.Models[i].Fields, id)
			}
			//	if ok{
			//  todo if the type of the field for id is not int .. it has to be int
			//	}
		}
	}

	return err
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
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

// WriteTmplToFile is to convert from template to a file
func WriteTmplToFile(filePath string, tmpl string, data interface{}) (err error) {
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
