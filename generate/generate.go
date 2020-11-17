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
	"strings"

	"golang.org/x/lint"
)

// Generater interface is to provide generater methods
type Generater interface {
	TmplToString(tmpl string, data interface{}) (result string, err error)
	TmplToFile(filePath string, tmpl string, data interface{}) (err error)
}

// Generate is a type
type Generate struct {
	Type       *string // Type of the project http , grpc , CloudEvents , cli
	Root       *string // ideally project root directory .i.e project name
	DBType     *string // mongo , sql based postgres mariadb etc
	HasHandler bool
	Models     []Model
	Gen        Generater
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
func New(file *string) (tg *Generate, err error) {
	if file == nil || *file == "" {
		return nil, errors.New("no file provided")
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

	root := strings.ToLower(*tg.Root)
	fmt.Println(root)
	*tg.Root = root

	err = tg.MkDirs() // Generate required directories
	if err != nil {
		return nil, err
	}

	return tg, nil
}

// MkDirs Create all required directories
func (tg *Generate) MkDirs() (err error) {
	if tg == nil || tg.Root == nil {
		return errors.New("project root directory or the generation has error")
	}
	err = os.Mkdir(*tg.Root, 0777)
	if err != nil {
		return err
	}

	models := filepath.Join(*tg.Root, "models")
	err = os.Mkdir(models, 0777)
	if err != nil {
		return err
	}

	if tg.HasHandler {
		handlers := filepath.Join(*tg.Root, "handlers")
		err = os.Mkdir(handlers, 0777)
		if err != nil {
			return err
		}
	}

	if tg.DBType != nil && *tg.DBType != "none" {
		database := filepath.Join(*tg.Root, "database")
		err = os.Mkdir(database, 0777)
		if err != nil {
			return err
		}
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

// CreateMain main.go
func (tg *Generate) CreateMain(tmpl string) (err error) {
	fileName := path.Join(*tg.Root, "main.go")

	data := make(map[string]string)

	data["project_name"] = *tg.Root

	if tg.Gen != nil {
		err = tg.Gen.TmplToFile(fileName, tmpl, data)
		if err != nil {
			return err
		}
	} else {
		err = TmplToFile(fileName, tmpl, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateAllModelFiles is to create all model files
func (tg *Generate) GenerateAllModelFiles(tmpl string) (err error) {
	for _, v := range tg.Models {
		modelsFile := path.Join(*tg.Root, "models", strings.ToLower(v.Name)+".go")
		if tg.Gen != nil {
			err = tg.Gen.TmplToFile(modelsFile, tmpl, v)
			if err != nil {
				return err
			}
		} else {
			err = TmplToFile(modelsFile, tmpl, v)
			if err != nil {
				return err
			}
		}

	}
	return err
}

// CopyAllStaticFiles is to copy static files
func (tg *Generate) CopyAllStaticFiles() (err error) {
	if tg == nil || tg.DBType == nil {
		return errors.New("Generate object or DBType is nil")
	}

	if *tg.DBType == "mongo" {
		database := filepath.Join(*tg.Root, "database", "mongo")
		err = os.Mkdir(database, 0777)
		if err != nil {
			return err
		}
	}

	src := filepath.Join("static", "databases", "mongo", "database.static")

	//src = strings.TrimSuffix(src, filepath.Ext(src)) + ".go"

	dst := filepath.Join(*tg.Root, "database", "mongo", "database.go")

	err = CopyStaticFile(src, dst)
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
