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
	"strings"

	"golang.org/x/lint"
)

// Generater interface is to provide generater methods
type Generater interface {
	ToString(tmpl string, data interface{}) (result string, err error)
	ToFile(filePath string, tmpl string, data interface{}) (err error)
	Read(key string) interface{}
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
	interfaces := filepath.Join(*tg.Root, "interfaces")
	err = os.Mkdir(interfaces, 0777)
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

	if tg.Type != nil && *tg.Type == "http" {
		database := filepath.Join(*tg.Root, "helper")
		err = os.Mkdir(database, 0777)
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

func (tg *Generate) GenerateAll() (err error) {
	if tg == nil {
		return errors.New("temp	late generater is not a valid object.Try to instantiate it through generater.New function")
	}
	if tg.Root == nil || *tg.Root == "" {
		return errors.New("invalid root directory")
	}
	// Todo write more conditions here

	ConfsMap := make(map[string][]string)
	ConfsMap["http_mongo"] = []string{"models", "interfaces", "database_mongo", "http_mongo_handler"}
	ConfsMap["only_models"] = []string{"models", "interfaces"}
	err = tg.CopyAllStaticFiles()
	if err != nil {
		return err
	}

	mhandler := make(map[string]interface{})
	mhandler["Root"] = tg.Root

	for _, v := range tg.Models {
		mhandler["Model"] = v

		modelsFile := path.Join(*tg.Root, "models", strings.ToLower(v.Name)+".go")
		err = GenerateFile(tg.Gen, mhandler, "models", modelsFile)
		if err != nil {
			return err
		}

		interfaceFile := path.Join(*tg.Root, "interfaces", strings.ToLower(v.Name)+".go")
		err = GenerateFile(tg.Gen, mhandler, "interfaces", interfaceFile)
		if err != nil {
			return err
		}

		daFile := path.Join(*tg.Root, "database", strings.ToLower(v.Name)+"DB.go")
		err = GenerateFile(tg.Gen, mhandler, "database_mongo", daFile)
		if err != nil {
			return err
		}
		handlerFile := path.Join(*tg.Root, "handlers", strings.ToLower(v.Name)+".go")
		err = GenerateFile(tg.Gen, mhandler, "http_mongo_handler", handlerFile)
		if err != nil {
			return err
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
	fileName := path.Join(*tg.Root, "main.go")

	data := make(map[string]string)

	data["project_name"] = *tg.Root

	if tg.Gen != nil {
		err = tg.Gen.ToFile(fileName, tmpl, data)
		if err != nil {
			return err
		}
	} else {
		errors.New("Gen interface is not assigned to Generate object")
	}
	return nil
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

// CopyAllStaticFiles is to copy static files
func (tg *Generate) CopyAllStaticFiles() (err error) {
	if tg == nil || tg.DBType == nil {
		return errors.New("Generate object or DBType is nil")
	}

	if *tg.DBType == "mongo" {

		src := filepath.Join("static", "databases", "mongo", "database.static")

		dst := filepath.Join(*tg.Root, "database", "database.go")

		err = CopyStaticFile(src, dst)
		if err != nil {
			return err
		}
	}

	if tg.Type != nil && *tg.Type == "http" {

		src := filepath.Join("static", "containers", "Dockerfile")

		dst := filepath.Join(*tg.Root, "Dockerfile")

		err = CopyStaticFile(src, dst)
		if err != nil {
			return err
		}

	}

	if tg.Type != nil && *tg.Type == "http" {

		src := filepath.Join("static", "helper", "helper.static")

		dst := filepath.Join(*tg.Root, "helper", "helper.go")

		err = CopyStaticFile(src, dst)
		if err != nil {
			return err
		}

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
