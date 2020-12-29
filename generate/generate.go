package generate

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"readyGo/mapping"
	"readyGo/scaler"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// New is to generate a new generater.
func New(file *string, mapping *mapping.Mapping, scaler scaler.Map) (tg *Generate, err error) {

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

	tg.Scalers = scaler

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
		}}).Parse(tmpl))
	err = t.Execute(file, data)

	if err != nil {
		return err
	}

	return nil
}
