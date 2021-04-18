package generate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrNoProjectName is to define error that project name is missing is not provided
	ErrNoProjectName = errors.New("project name is missing")
	//ErrInvalidProjectType is to define error as invalid project type
	ErrInvalidProjectType = errors.New("invalid project type;project type must be http | grpc | cloudEvent | cli")
	// ErrInvalidDatabase is to define error that invalid database
	ErrInvalidDatabase          = errors.New("invalid database;database type (DB) must be mongo | sql ")
	ErrInvalidDatabaseName      = errors.New("sql kind supports only mysql or postgress")
	ErrInvalidDatabasenosqlName = errors.New("nosql kind supports only mongo or cassendra")
	ErrInvalidDatabaseKind      = errors.New("database kind supports only sql or nosql")
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
		if !tg.Implementer.IsValidIdentifier(m.Name) {
			fmt.Println(m.Name)
			return errors.New(m.Name + " is invalid identifier")
		}
		tmpModel := m.Name
		tg.Models[i].Name = strings.ToUpper(string(tmpModel[0])) + string(tmpModel[1:])
		for j, f := range m.Fields {
			if !tg.Implementer.IsValidIdentifier(f.Name) {
				fmt.Println(f.Name)
				return errors.New(f.Name + " is invalid identifier")
			}
			tmpField := f.Name
			tg.Models[i].Fields[j].Name = strings.ToUpper(string(tmpField[0])) + string(tmpField[1:])
		}
	}
	return nil
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
	if tg.DatabaseSpec.Kind == "sql" && (tg.DatabaseSpec.Name != "mysql" && tg.DatabaseSpec.Name != "postgres") {
		return ErrInvalidDatabaseName
	}
	if tg.DatabaseSpec.Kind == "nosql" && (tg.DatabaseSpec.Name != "mongo" && tg.DatabaseSpec.Name != "cassendra") {
		return ErrInvalidDatabasenosqlName
	}
	if tg.DatabaseSpec.Kind != "nosql" && tg.DatabaseSpec.Kind != "sql" {
		return ErrInvalidDatabaseKind
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
		if tg.DatabaseSpec.Name == "mongo" {
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

// SetFieldCategory each field is set into certain category. They are scaler, model , function,default value etc..
func (tg *Generate) SetFieldCategory() error {
	for mi := 0; mi < len(tg.Models); mi++ {
		for i, f := range tg.Models[mi].Fields {
			ftype := f.Type
			if strings.Contains(f.Type, "[]") {
				ftype = strings.Split(ftype, "[]")[1]
			}
			if tg.Scalers.IsValidreadyGotype(ftype) {
				tg.Models[mi].Fields[i].Category = "scaler" // all readyGo types are scaler types
			} else if strings.Contains(ftype, "global.") {
				tg.Models[mi].Fields[i].Definition = ftype
				tg.Models[mi].Fields[i].Type = tg.Implementer.GetFuncReturnType(ftype) // Todo fetch this from reading global. functions return type
				tg.Models[mi].Fields[i].Category = "function"                          // any type that contains global. is a function category as all functionsa re global.xxxfuncname
			} else if tg.IsModelType(ftype) {
				fmt.Println(ftype, f.Type)
				if strings.Contains(f.Type, "[]") {
					tg.Models[mi].Fields[i].Category = "array model" // model category are types that are already one of the models
				} else {
					tg.Models[mi].Fields[i].Category = "model" // model category are types that are already one of the models
				}
			} else {
				tg.Models[mi].Fields[i].Category = "undefined" // undefined category is a category that does not fall into any of above
				return errors.New(ftype + " type is undefined")
			}
		}
	}
	return nil
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

// IsModelArrayType is to check whether a filed is a model field
func (tg *Generate) IsModelArrayType(iden string) bool {
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

// ValidateTypes checks whether a type is either a readyGo or model type
func (tg *Generate) ValidateTypes() error {
	//Use Go's shorthand range loop
	for m := range tg.Models {
		for f := range tg.Models[m].Fields {
			// Check model name equals field type first, preventing unnecessary map lookups
			if tg.Models[m].Name == tg.Models[m].Fields[f].Type {
				continue
			}

			// Remove function IsValidreadyGotype() overhead. It's not needed for a simple map lookup.
			// Function names should have correct capitalisation. I would have renamed to IsReadyGoType() but the function isn't needed.
			if _, ok := tg.Scalers[tg.Models[m].Fields[f].Type]; ok {
				continue
			}

			var found bool
			for md := range tg.Models {
				// Ignore the current model because it has already been checked.
				if md == m {
					continue
				}

				// Use indexes to lookup rather than copying the struct to another variable.
				if tg.Models[md].Name == tg.Models[m].Fields[f].Type {
					found = true
					// Exit out of the loop as soon as a match is found.
					break
				}
			}
			if !found {
				// Give more detail in error messages so you spend less time debugging. E.g: What if the type was misspelled as a tab, null (byte(0)) or space characters?
				return fmt.Errorf("`%s` is neither a readyGo type nor a model type in `%s`", tg.Models[m].Fields[f].Type, tg.Models[m].Name)
			}
		}
	}

	return nil
}
