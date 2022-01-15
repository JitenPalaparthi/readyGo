package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
	"readyGo/box"
	"readyGo/helper"
	"readyGo/mapping"
	"readyGo/scalar"
)

var (
	// ErrInvalidProjectName is to define error information upon invalid project name
	ErrInvalidProjectName = errors.New("invalid project name; it must have only characters;no special chars,whitespaces,digits are allowed;")

	// ErrNoFile is to define error that no file provided
	ErrNoFile = errors.New("no file provided")

	// ErrEmptyMapping is to define error that mapping is empty.
	ErrEmptyMapping = errors.New("invalid mapping;mapping cannot be empty")

	// ErrEmptyImplementer is to define error that implementer is nil
	ErrEmptyImplementer = errors.New("invalid implementer;implmenter cannot be nil")

	// ErrInvalidTemplateGenerator is to define error that invalid template generator is provided
	ErrInvalidTemplateGenerator = errors.New("invalid template generator;try to instantiate it through generator.New function")

	// ErrInvalidRoot is to define error as invalid root directory
	ErrInvalidRoot = errors.New("invalid root directory")
)

// New is to generate a new generater.
func New(file *string, scalar scalar.Map, implementer Implementer) (tg *Generate, err error) {

	if file == nil || *file == "" {
		return nil, ErrNoFile
	}

	if implementer == nil {
		return nil, ErrEmptyImplementer
	}

	ext := filepath.Ext(*file)
	if ext != ".json" && ext != ".yaml" && ext != ".yml" {
		return nil, errors.New("only json | yaml | yml files are allowed ")
	}
	cFile, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, err
	}

	switch ext {
	case ".json":
		err = json.Unmarshal(cFile, &tg)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(cFile, &tg)
	}
	if err != nil {
		return nil, err
	}

	// Logic to identify mapping
	projectType := tg.APISpec.Kind + "_" + tg.DatabaseSpec.Name
	if tg.MessagingSpec.Name != "" {
		projectType = projectType + "_" + tg.MessagingSpec.Name
	}

	//ops := boxops.New("../box")
	ops := &box.Box{}
	mapping, err := mapping.New(ops, filepath.Join("configs", "mappings", projectType+".json"), projectType)
	if err != nil {
		log.Fatal(err)
	}

	tg.Mapping = mapping

	tg.Scalars = scalar

	tg.Implementer = implementer

	err = tg.ChangeIden()
	if err != nil {
		return nil, err
	}

	matched, err := regexp.MatchString("^[a-zA-Z]*$", tg.Project)
	if !matched {
		return nil, ErrInvalidProjectName
	}
	if err != nil {
		return nil, err
	}

	err = tg.SetFieldCategory()
	if err != nil {
		return tg, err
	}

	err = tg.Validate()

	if err != nil {
		return nil, err
	}

	// This channel can be used in such a way that all generated output can be sent to this so that I can be prined properly
	if tg.Output == nil {
		tg.Output = make(chan string)
	}

	return tg, nil
}

// CreateAll creates all kinds of files based on the provided mappings
func (tg *Generate) CreateAll() (err error) {
	if tg == nil {
		return ErrInvalidTemplateGenerator
	}
	if tg.Project == "" {
		return ErrInvalidRoot
	}
	// Todo write more conditions here

	for _, opsData := range tg.Mapping.OpsData {
		//fmt.Println(opsData)
		switch opsData.OpType {
		case "directories":
			path := filepath.Join(tg.Project, opsData.Src)
			tg.Output <- "generating the following directory :" + path
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return errorsJoin(err, tg.RmDir())
			}
			tg.Output <- "the following directory has been generated :" + path
		case "static-files":
			dst := filepath.Join(tg.Project, opsData.Dst)
			tg.Output <- "generating the following static file :" + dst
			content, err := tg.Mapping.Reader.Read(opsData.Src)
			if err != nil {
				return errorsJoin(err, tg.RmDir())
			}
			if content != "" {
				var li int
				switch runtime.GOOS {
				case "windows":
					li = strings.LastIndex(dst, "\\")
				default:
					li = strings.LastIndex(dst, "/")
				}

				dirs := dst[0:li]
				err = os.MkdirAll(dirs, 0755)
				if err != nil {
					return errorsJoin(err, tg.RmDir())
				}

				err = ioutil.WriteFile(dst, []byte(content), 0644)
				if err != nil {
					return errorsJoin(err, tg.RmDir())
				}
				tg.Output <- "The following static file has been generated :" + dst
			}
		case "multiple-file-templates":
			mHandler := map[string]interface{}{
				"Project": tg.Project,
				"config":  tg,
			}
			if opsData.GenForType == "both" {
				for _, v := range tg.Models {
					mHandler["Model"] = v
					dst := filepath.Join(tg.Project, opsData.Dst)
					tg.Output <- "generating template based file :" + dst
					err = os.MkdirAll(dst, 0755)
					if err != nil {
						return errorsJoin(err, tg.RmDir())
					}
					// If there is any extension in the opsData that means file to be created with the given extension. Otherwise create a default one with .go
					if opsData.Ext == "" {
						dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+".go")
					} else {
						// If any extension starts with . add the extension as it is.Otherwise add . as a prefix to the opsData.Ext
						if string(strings.TrimSpace(opsData.Ext)[0]) == "." {
							dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+opsData.Ext)
						} else {
							dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+"."+opsData.Ext)
						}
					}
					content, err := tg.Mapping.Reader.Read(opsData.Src)
					if err != nil {
						return errorsJoin(err, tg.RmDir())
					}
					if content != "" {
						err = tg.WriteTmplToFile(dst, content, mHandler)
						if err != nil {
							return errorsJoin(err, tg.RmDir())
						}
						tg.Output <- "the following templated based file has been generated :" + dst
					}
				}
			} else if opsData.GenForType == "main" {
				for _, v := range tg.Models {
					if v.Type == "main" {
						err = generateTemplateFile(&mHandler, v, tg, opsData)
						if err != nil {
							return err
						}
					}
				}
			} else if opsData.GenForType == "sub" {
				for _, v := range tg.Models {
					if v.Type == "sub" {
						err = generateTemplateFile(&mHandler, v, tg, opsData)
						if err != nil {
							return err
						}
					}
				}
			}

		case "single-file-templates":
			_, _, err = generateFiles(tg, opsData, "template", "generating template based contents to the file :")
			if err != nil {
				return err
			}

		case "exec":
			var content, dst string
			content, dst, err = generateFiles(tg, opsData, "shell", "generating shell based executable files :")
			if err != nil {
				return err
			}

			if content != "" {
				err = os.Chmod(dst, 0700)
				if err != nil {
					return err
				}
				tg.Output <- "giving read|writeexecute permissions to the file :" + dst
			}
		default:
			return errors.New(opsData.OpType + ":this type has no implementation")
		}
	}

	return nil
}

func generateFiles(tg *Generate, opsData mapping.OpsData, fileType, output string) (content, dst string, err error) {
	// Todo for opsData.Ext if there is an extension
	dst = path.Join(tg.Project, opsData.Dst)
	tg.Output <- output + dst

	li := strings.LastIndex(dst, "/")
	dirs := dst[:li]

	err = os.MkdirAll(dirs, 0755)
	if err != nil {
		return "", dst, errorsJoin(err, tg.RmDir())
	}

	content, err = tg.Mapping.Reader.Read(opsData.Src)
	if err != nil {
		return "", dst, errorsJoin(err, tg.RmDir())
	}

	if content != "" {
		mHandler := map[string]interface{}{
			"config": tg,
		}
		err = tg.WriteTmplToFile(dst, content, mHandler)
		if err != nil {
			return "", dst, errorsJoin(err, tg.RmDir())
		}

		tg.Output <- fmt.Sprintf("The following %s file has been generated :%s", fileType, dst)
	}

	return content, dst, nil
}

// RmDir is to remove dirs
func (tg *Generate) RmDir() (err error) {
	tg.Output <- "removing all directories of the project :" + tg.Project
	err = os.RemoveAll(tg.Project)
	if err != nil {
		return err
	}
	tg.Output <- "all directories in the following project has been removed :" + tg.Project
	return nil
}

// WriteTmplToFile is to convert from template to a file
func (tg *Generate) WriteTmplToFile(filePath string, tmpl string, data interface{}) (err error) {
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
		}).Funcs(template.FuncMap{
		"Counter": func(str string) string {
			if count, err := strconv.Atoi(str); err == nil {
				return strconv.Itoa(count + 1)
			}
			return "0"
		}}).Funcs(template.FuncMap{
		"GoType": func(tpe string) string {
			if scler := tg.Scalars.GetScalar(tpe); scler != nil {
				return scler.GoType
			}
			return ""
		}}).Funcs(template.FuncMap{
		"GrpcType": func(tpe string) string {
			if scler := tg.Scalars.GetScalar(tpe); scler != nil {
				return scler.GrpcType
			}
			return ""
		}}).Funcs(template.FuncMap{
		"GrpcArrayModel": func(tpe string) string {
			ss := strings.Split(tpe, "[]")
			if len(ss) > 1 {
				return ss[1]
			}
			return ""
		}}).Funcs(template.FuncMap{
		"GoRegExFormat": func(str string) string {
			if str == "" {
				return ""
			}
			str = strings.TrimSpace(str)
			//strbuff := []byte(str)
			if len(str) > 2 {
				//	strbuff[0] = 96
				//	strbuff[len(strbuff)-1] = 96
				stroriginal := str
				str = strings.Replace(str[1:len(str)-1], "`", `"`+"`"+`"`, -2)
				return string(stroriginal[0]) + str + string(stroriginal[len(stroriginal)-1])
			}
			return str
		}}).Parse(tmpl))

	return t.Execute(file, data)
}

func ExecuteCommand(filename string) (string, error) {
	cmd, err := exec.Command("/bin/sh", filename).Output()
	if err != nil {
		return "", err
	}

	return string(cmd), nil
}

// Execute executes given shell files
func (tg *Generate) Execute() (err error) {
	if !helper.IsWindows() {
		if tg == nil {
			return ErrInvalidTemplateGenerator
		}
		if tg.Project == "" {
			return ErrInvalidRoot
		}
		for _, opsData := range tg.Mapping.OpsData {
			switch opsData.OpType {
			case "exec":
				// Todo for opsData.Ext if there is an extension
				tg.Output <- "executing the following file:" + opsData.Dst
				mhandler := make(map[string]interface{})
				mhandler["config"] = tg
				dst := path.Join(tg.Project, opsData.Dst)
				output, err := ExecuteCommand(dst)
				if err != nil {
					return err
				}
				tg.Output <- "the following file has been executed:" + opsData.Dst
				tg.Output <- output // Sending output to the channel
			default:
			}
		}
	}
	return nil
}

// WriteOutput this should be a go routine
func (tg *Generate) WriteOutput(w io.Writer) {
	for output := range tg.Output {
		output = "\n" + output // Add a new line to the output
		_, err := w.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateTemplateFile(mHandler *map[string]interface{}, v Model, tg *Generate, opsData mapping.OpsData) (err error) {
	(*mHandler)["Model"] = v
	dst := filepath.Join(tg.Project, opsData.Dst)
	tg.Output <- "generating template based file :" + dst

	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return errorsJoin(err, tg.RmDir())
	}

	// If there is any extension in the opsData that means file to be created with the given extension. Otherwise create a default one with .go
	if opsData.Ext == "" {
		dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+".go")
	} else {
		// If any extension starts with . add the extension as it is.Otherwise add . as a prefix to the opsData.Ext
		if string(strings.TrimSpace(opsData.Ext)[0]) == "." {
			dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+opsData.Ext)
		} else {
			dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+"."+opsData.Ext)
		}
	}

	content, err := tg.Mapping.Reader.Read(opsData.Src)
	if err != nil {
		return errorsJoin(err, tg.RmDir())
	}
	if content != "" {
		err = tg.WriteTmplToFile(dst, content, *mHandler)
		if err != nil {
			return errorsJoin(err, tg.RmDir())
		}
		tg.Output <- "the following templated based file has been generated :" + dst
	}

	return err
}

func errorsJoin(err1, err2 error) error {
	if err2 != nil {
		return errors.New("1." + err1.Error() + ".2." + err2.Error())
	}

	return err1
}
