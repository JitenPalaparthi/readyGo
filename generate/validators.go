package generate

import (
	"errors"
	"log"
	"sort"
	"strings"

	"golang.org/x/lint"
)

// ValidateAndChangeIdentifier is to validate and Change as and where required
func (tg *Generate) ValidateAndChangeIdentifier() (err error) {
	for i, m := range tg.Models {
		tmpModel := m.Name
		ok, problems, err := ValidModelType(tmpModel)
		if !ok {
			for _, problem := range problems {
				log.Println("Warning:", problem)
			}
			if err != nil {
				return err
			}
		}
		tg.Models[i].Name = strings.ToUpper(string(tmpModel[0])) + string(tmpModel[1:])
		for j, f := range m.Fields {
			tmpField := f.Name
			ok, problems, err := ValidFieldType(tmpField, f.Type)
			if !ok {
				for _, problem := range problems {
					log.Println("Warning:", problem)
				}
				if err != nil {
					return err
				}
				//return errors.New("invalid field or type in configuration file :" + tmpField + ":" + f.Type)
			}
			tg.Models[i].Fields[j].Name = strings.ToUpper(string(tmpField[0])) + string(tmpField[1:])
		}
	}
	return nil

}

// ValidModelType is to validate the model type
func ValidModelType(name string) (bool, []lint.Problem, error) {
	var l lint.Linter
	problems, err := l.Lint("", []byte("//Package valid is a valid package \n package valid \n type "+name+" struct{}"))
	//fmt.Println(problems, err)
	if len(problems) > 0 || err != nil {
		return false, problems, err
	}
	return true, nil, nil
}

// ValidFieldType is to validate the model type
func ValidFieldType(field, tpe string) (bool, []lint.Problem, error) {
	var l lint.Linter
	problems, err := l.Lint("", []byte("//Package valid is a valid package \n package valid\n type "+field+" "+tpe))
	//	fmt.Println(problems, err)
	if len(problems) > 0 || err != nil {
		return false, problems, err
	}
	return true, nil, nil
}

// Validate is to validate the object
func (tg *Generate) Validate() (err error) {
	if tg.Project == "" {
		return errors.New("Project name is missing")
	}

	// Validate whether Project name is proper identifier

	if tg.Kind == "" || (tg.Kind != "http" && tg.Kind != "grpc" && tg.Kind != "cloudEvent" && tg.Kind != "cli") {
		return errors.New(" Project type must be http | grpc | cloudEvent | cli")
	}
	if tg.DatabaseSpec.Kind == "" || (tg.DatabaseSpec.Kind != "mongo" && tg.DatabaseSpec.Kind != "sql") {
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
		if tg.DatabaseSpec.Kind == "mongo" {
			_, ok := fieldMap["id"]
			if !ok {
				id := Field{Name: "Id", Type: "string"}
				tg.Models[i].Fields = append(tg.Models[i].Fields, id)
			}
			//	if ok{
			//  todo if the type of the field for id is not string .. it has to be string
			//	}
		}
		if tg.DatabaseSpec.Kind == "sql" {
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
