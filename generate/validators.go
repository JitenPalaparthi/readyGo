package generate

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/lint"
)

var (
	// ErrNoProjectName is to define error that project name is missing is not provided
	ErrNoProjectName = errors.New("project name is missing")
	//ErrInvalidProjectType is to define error as invalid project type
	ErrInvalidProjectType = errors.New("invalid project type;project type must be http | grpc | cloudEvent | cli")
	// ErrInvalidDatabase is to define error that invalid database
	ErrInvalidDatabase = errors.New("invalid database;databas type (DB) must be mongo | sql ")
)

// IsValidIdentifier is to check whether the field is a valid identifier or not
func (tg *Generate) IsValidIdentifier(fielden string) bool {

	// Should not start with the number or any special chars other than _
	// should not contain secial chars other than _
	// should have atleast one char and can have n number of digits

	// according to https://www.geeksforgeeks.org/check-whether-the-given-string-is-a-valid-identifier/
	// It must start with either underscore(_) or any of the characters from the ranges [‘a’, ‘z’] and [‘A’, ‘Z’].
	// There must not be any white space in the string.
	//  And, all the subsequent characters after the first character must not consist of any special characters like $, #, % etc.

	// This tests whether a pattern matches a string.
	match, err := regexp.MatchString(`^[^\d\W]\w*$`, fielden)
	if err != nil {
		return false
	}
	return match
}

// ChangeIden is to change all identifiers based on
func (tg *Generate) ChangeIden() error {
	for i, m := range tg.Models {
		if !tg.IsValidIdentifier(m.Name) {
			fmt.Println(m.Name)
			return errors.New(m.Name + " is invalid identifier")
		}
		tmpModel := m.Name
		tg.Models[i].Name = strings.ToUpper(string(tmpModel[0])) + string(tmpModel[1:])
		for j, f := range m.Fields {
			if !tg.IsValidIdentifier(f.Name) {
				fmt.Println(f.Name)
				return errors.New(f.Name + " is invalid identifier")
			}
			tmpField := f.Name
			tg.Models[i].Fields[j].Name = strings.ToUpper(string(tmpField[0])) + string(tmpField[1:])
		}
	}
	return nil
}

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
		return ErrNoProjectName
	}
	// Validate whether Project name is proper identifier

	if tg.Kind == "" || (tg.Kind != "http" && tg.Kind != "grpc" && tg.Kind != "cloudEvent" && tg.Kind != "cli") {
		return ErrInvalidProjectType
	}
	if tg.DatabaseSpec.Kind == "" || (tg.DatabaseSpec.Kind != "mongo" && tg.DatabaseSpec.Kind != "sql") {
		return ErrInvalidDatabase
	}
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

		}
		if tg.DatabaseSpec.Kind == "mongo" {
			_, ok := fieldMap["id"]
			if !ok {
				id := Field{Name: "ID", Type: "string", Category: "scaler"}
				tg.Models[i].Fields = append(tg.Models[i].Fields, id)
			}
			//	if ok{
			//  todo if the type of the field for id is not string .. it has to be string
			//	}
		}
		if tg.DatabaseSpec.Kind == "sql" {
			_, ok := fieldMap["id"]
			if !ok {
				id := Field{Name: "ID", Type: "int", Category: "scaler"}
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

// SetFieldCategory each field is set into certain category. They are scaler, model , function,default value etc..
func (tg *Generate) SetFieldCategory() {
	for mi := 0; mi < len(tg.Models); mi++ {
		for i, f := range tg.Models[mi].Fields {
			if tg.Scalers.IsValidreadyGotype(f.Type) {
				tg.Models[mi].Fields[i].Category = "scaler" // all readyGo types are scaler types
			} else if strings.Contains(f.Type, "global.") {
				tg.Models[mi].Fields[i].Category = "function" // any type that contains global. is a function category as all functionsa re global.xxxfuncname
			} else if tg.IsModelType(f.Type) {
				tg.Models[mi].Fields[i].Category = "model" // model category are types that are already one of the models
			} else {
				tg.Models[mi].Fields[i].Category = "undefined" // undefined category is a category that does not fall into any of above
			}
		}
	}
}

// IsModelType is to check whether a filed is a model field
func (tg *Generate) IsModelType(iden string) bool {
	models := make([]string, 0)
	for mi := 0; mi < len(tg.Models); mi++ {
		models = append(models, tg.Models[mi].Name)
	}
	for _, m := range models {
		if m == iden {
			return true
		}
	}
	return false
}

// ValidateTypes is to valudate whether a type is readyGo type or a model type
func (tg *Generate) ValidateTypes() (err error) {
	models := make([]string, 0)
	fields := make([]string, 0)
	for mi := 0; mi < len(tg.Models); mi++ {
		models = append(models, tg.Models[mi].Name)
		for _, f := range tg.Models[mi].Fields {
			if !tg.Scalers.IsValidreadyGotype(f.Type) {
				fields = append(fields, f.Type)
			}
		}
	}
	if len(fields) == 0 {
		return nil
	}
	var check bool = false

	for _, f := range fields {
		for _, md := range models {
			if md == f {
				check = true
			}
		}
		if !check {
			return errors.New(f + " is neither a readyGo type nor a model type")
		}
	}
	return nil
}
