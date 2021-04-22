package generate

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"readyGo/boxops"
	"readyGo/mapping"
	"readyGo/scaler"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"readyGo/helper"

	"gopkg.in/yaml.v2"
)

var (
	//ErrInvalidProjectName is to define error information upon invalid project name
	ErrInvalidProjectName = errors.New("invalid project name; it must have only characters;no special chars,whitespaces,digits are allowed;")

	// ErrNoFile is to define error that no file provided
	ErrNoFile = errors.New("no file provided")

	// ErrEmptyMapping is to define error that mapping is empty.
	ErrEmptyMapping = errors.New("invalid mapping;mapping cannot be empty")

	//ErrEmptyImplementer is to define error that implementer is nil
	ErrEmptyImplementer = errors.New("invalid implementer;implmenter cannot be nil")

	// ErrInvalidTemlateGenerator is to define error that invalid template generator is provided
	ErrInvalidTemlateGenerator = errors.New("invalid template generater;try to instantiate it through generater.New function")

	// ErrInvalidRoot is to define error as invalid root directory
	ErrInvalidRoot = errors.New("invalid root directory")
)

// New is to generate a new generater.
func New(file *string, scaler scaler.Map, implementer Implementer) (tg *Generate, err error) {

	if file == nil || *file == "" {
		return nil, ErrNoFile
	}

	if implementer == nil {
		return nil, ErrEmptyImplementer
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

	// Logic to identify mapping

	projectType := tg.APISpec.Kind + "_" + tg.DatabaseSpec.Name
	if tg.MessagingSpec.Name != "" {
		projectType = projectType + "_" + tg.MessagingSpec.Name
	}

	ops := boxops.New("../box")
	mapping, err := mapping.New(ops, "configs/mappings/"+projectType+".json", projectType)
	if err != nil {
		log.Fatal(err)
	}

	tg.Mapping = mapping

	tg.Scalers = scaler

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
		return ErrInvalidTemlateGenerator
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
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				errRm := tg.RmDir()
				if errRm != nil {
					return errors.New("1." + err.Error() + ".2." + errRm.Error())
				}
				return err
			}
			tg.Output <- "creating directory :" + path
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
				tg.Output <- "writing contents to the file :" + dst
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
				// If there is any extension in the opsData that means file to be created with the given extension. Otherwise create a default one with .go
				if opsData.Ext == "" {
					dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+".go")
				} else {
					// If any extension starts with . add the extension as it is.Otherwise add . as a prefix to the opsData.Ext
					if string(strings.Trim(opsData.Ext, " ")[0]) == "." {
						dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+opsData.Ext)
					} else {
						dst = path.Join(tg.Project, opsData.Dst, strings.ToLower(v.Name)+"."+opsData.Ext)
					}
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
					err := tg.WriteTmplToFile(dst, content, mhandler)
					if err != nil {
						errRm := tg.RmDir()
						if errRm != nil {
							return errors.New("1." + err.Error() + ".2." + errRm.Error())
						}
						return err
					}
					tg.Output <- "writing template based contents to the file :" + dst
				}
			}

		case "single-file-templates":
			// Todo for opsData.Ext if there is an extension
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
				err := tg.WriteTmplToFile(dst, content, mhandler)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}
				tg.Output <- "writing template based contents to the file :" + dst
			}
		case "exec":
			// Todo for opsData.Ext if there is an extension
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
				err := tg.WriteTmplToFile(dst, content, mhandler)
				if err != nil {
					errRm := tg.RmDir()
					if errRm != nil {
						return errors.New("1." + err.Error() + ".2." + errRm.Error())
					}
					return err
				}
				tg.Output <- "writing shall based executable files :" + dst

				os.Chmod(dst, 0700)
				tg.Output <- "giving read|writeexecute permissions to the file :" + dst

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
	tg.Output <- "removing all directories  of the project :" + tg.Project
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
			if s, err := strconv.Atoi(str); err == nil {
				count := s + 1
				return strconv.Itoa(count)
			}
			return "0"
		}}).Funcs(template.FuncMap{
		"GoType": func(tpe string) string {
			if scler := tg.Scalers.GetScaler(tpe); scler != nil {
				return scler.GoType
			}
			return ""
		}}).Funcs(template.FuncMap{
		"GrpcType": func(tpe string) string {
			if scler := tg.Scalers.GetScaler(tpe); scler != nil {
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
			str = strings.Trim(str, " ")
			//strbuff := []byte(str)
			if len(str) > 2 {
				//	strbuff[0] = 96
				//	strbuff[len(strbuff)-1] = 96
				stroriginal := str
				str = strings.Replace(str[1:len(str)-1], "`", `"`+"`"+`"`, -2)
				return string(stroriginal[0]) + str + string(stroriginal[len(stroriginal)-1])
			}
			return string(str)
		}}).Parse(tmpl))
	err = t.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}

func ExecuteCommand(filename string) (string, error) {
	cmd, err := exec.Command("/bin/sh", filename).Output()
	if err != nil {
		return "", err
	}
	output := string(cmd)
	return output, nil
}

// Execute executes given shell files
func (tg *Generate) Execute() (err error) {
	if helper.IsWindows() {
		return errors.New("exec feature works only for unix based OS")
	}
	if tg == nil {
		return ErrInvalidTemlateGenerator
	}
	if tg.Project == "" {
		return ErrInvalidRoot
	}
	for _, opsData := range tg.Mapping.OpsData {
		switch opsData.OpType {
		case "exec":
			// Todo for opsData.Ext if there is an extension
			mhandler := make(map[string]interface{})
			mhandler["config"] = tg
			dst := path.Join(tg.Project, opsData.Dst)
			output, err := ExecuteCommand(dst)
			if err != nil {
				return err
			}
			tg.Output <- "executing the file:" + opsData.Dst
			tg.Output <- output // Sending output to the channel
		default:
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
