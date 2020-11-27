package templates

import (
	template "readyGo/generate/template"
)

func New() template.TmplMap {
	var tm template.TmplMap
	tm = make(map[string]string)

	// config constatnt
	tm["config"] = `{
    "connection":"mongodb://local:27017",
    "db":"{{.Project}}"
}`

	tm["database"] = `package database

import (
	"context"
	"errors"
	"{{.Project}}/helper"
	"{{.Project}}/models"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// {{.Model.Name}}DB is to maintain database related methods
type {{.Model.Name}}DB struct {
	DB *Database
}

{{- $modelName:=.Model.Name}}
{{- range .Model.Fields}}
{{- if eq .IsKey true }}
func ({{$modelName | Initial}} *{{$modelName}}DB) Is{{$modelName}}ExistsBy{{.Name}}({{.Name | ToLower}} {{.Type}}) bool {
	if {{.Name | ToLower}} == "" {
		return false
	}
	filter := make(map[string]interface{}, 0)
	filter["{{.Name | ToLower}}"] = {{.Name | ToLower}}
	count, err := {{$modelName | Initial}}.DB.GetCount("{{$modelName | ToLower}}s", filter)
	if err != nil {
		if err.Error() == "not found" {
			return false
		}
	}
	if count > 0 {
		return true
	}
	return false
}
{{- end}}
{{- end}}

func ({{.Model.Name | Initial}} *{{.Model.Name}}DB) Create{{.Model.Name}}({{.Model.Name | ToLower}} *models.{{.Model.Name}}) (result interface{}, err error) {
		{{- $modelName:=.Model.Name}}
		{{- range .Model.Fields}}
		{{- if eq .IsKey true }}
		if {{$modelName | Initial}}.Is{{$modelName}}ExistsBy{{.Name}}({{$modelName | ToLower}}.{{.Name}}) {
		return nil,errors.New("{{$modelName}} already existed")
		}
		{{- end}}
		{{- end}}
		data, err := helper.ToMap({{.Model.Name | ToLower}}, "bson", "_id,omitempty")
		if err != nil {
			return nil, err
		}
		result, err = {{.Model.Name | Initial}}.DB.InsertRecord("{{.Model.Name | ToLower}}s", data)
		if err != nil {
			return nil, err
		}
	return result, nil
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}DB) Update{{.Model.Name}}ByID(id string, data map[string]interface{}) (result interface{}, err error) {
	result, err = {{.Model.Name | Initial}}.DB.UpdateRecordByID("{{.Model.Name | ToLower}}s", id, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}DB) Delete{{.Model.Name}}ByID(id string) (result interface{}, err error) {
	result, err = {{.Model.Name | Initial}}.DB.DeleteRecordByID("{{.Model.Name | ToLower}}s", id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}DB) Get{{.Model.Name}}ByID(id string) (*models.{{.Model.Name}}, error) {
	{{.Model.Name | ToLower}} := &models.{{.Model.Name}}{}
	mapData, err := {{.Model.Name | Initial}}.DB.FindRecordByID("{{.Model.Name | ToLower}}s",id)
	if err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(mapData, &{{.Model.Name | ToLower}}); err != nil {
		return nil, err
	}
	_id := mapData.(map[string]interface{})["_id"].(primitive.ObjectID).Hex()
	{{.Model.Name | ToLower}}.Id = _id
	return {{.Model.Name | ToLower}}, nil
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}DB) GetAll{{.Model.Name}}s(skip int64, limit int64, selector interface{}) ([]models.{{.Model.Name}}, error) {
	if _, ok := selector.(map[string]interface{}); !ok {
		return nil, errors.New("wrong input type")
	}
	var result []models.{{.Model.Name}}
	colleection := {{.Model.Name | Initial}}.DB.Client.(*mongo.Client).Database({{.Model.Name | Initial}}.DB.Name).Collection("{{.Model.Name | ToLower}}s")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	findOptions := options.Find()
	findOptions.SetLimit(limit).SetSkip(skip)
	cur, err := colleection.Find(ctx, selector, findOptions)
	if err != nil {
		return nil, err
	}
	result = make([]models.{{.Model.Name}}, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		iresult := models.{{.Model.Name}}{}
		err := cur.Decode(&iresult)
		if err != nil {
			return nil, err
		}
		result = append(result, iresult)
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}DB) GetAll{{.Model.Name}}sBy(search string, selector interface{}, skip int64, limit int64) ([]models.{{.Model.Name}}, error) {
	if _, ok := selector.(map[string]interface{}); !ok {
		return nil, errors.New("wrong input type")
	}
	var result []models.{{.Model.Name}}
	colleection := {{.Model.Name | Initial}}.DB.Client.(*mongo.Client).Database({{.Model.Name | Initial}}.DB.Name).Collection("{{.Model.Name | ToLower}}s")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	findOptions := options.Find()
	findOptions.SetLimit(limit).SetSkip(skip)

	if search != "" {
		selector.(map[string]interface{})["$text"] = bson.M{"$search": search}
	}

	cur, err := colleection.Find(ctx, selector, findOptions)
	if err != nil {
		return nil, err
	}
	result = make([]models.{{.Model.Name}}, 0)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		iresult := models.{{.Model.Name}}{}
		err := cur.Decode(&iresult)
		if err != nil {
			return nil, err
		}
		result = append(result, iresult)
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
`
	tm["docker-compose"] = `version: "3"
services:
  app:
    container_name: {{$.config.Project}}_service
    restart: always
    build: ./
    ports:
     - "{{$.config.Port}}:{{$.config.Port}}"
    links:
       - mongo
    networks: 
      - backend
  
  mongo:
    container_name: mongo_db
    image: mongo
    ports:
     - "27017:27017"
    volumes:
     - "./data/db:/data/db"
    networks: 
      - backend

networks:
  backend:
    driver: bridge
`
	tm["handlers"] = `// Author : readyGo "JitenP@Outlook.Com"
// This code is generated by template engine. You are free to make amendments as and where required
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"{{.Project}}/interfaces"
	"{{.Project}}/models"

	"github.com/gin-gonic/gin"
)

type {{.Model.Name}} struct {
	I{{.Model.Name}} interfaces.{{.Model.Name}}Interface
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}) Create{{.Model.Name}}() func(c *gin.Context) {
	var err error
	return func(c *gin.Context) {
		if c.Request.Method == "POST" {
			var {{.Model.Name | ToLower}} *models.{{.Model.Name}}
			{{.Model.Name | ToLower}} = &models.{{.Model.Name}}{}

			err = json.NewDecoder(c.Request.Body).Decode(&{{.Model.Name | ToLower}})
			fmt.Println({{.Model.Name | ToLower}})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			// Validate model
			err = models.Validate{{.Model.Name}}({{.Model.Name | ToLower}})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			result, err := {{.Model.Name | Initial}}.I{{.Model.Name}}.Create{{.Model.Name}}({{.Model.Name | ToLower}})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": result,
			})
			c.Abort()
			return
		}
	}
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}) Get{{.Model.Name}}ByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			id := c.Param("id")
			if id == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "id parameter has not been provided",
				})
				c.Abort()
				return
			}
			{{.Model.Name | ToLower}}, err := {{.Model.Name | Initial}}.I{{.Model.Name}}.Get{{.Model.Name}}ByID(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, {{.Model.Name | ToLower}})
		}
	}
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}) GetAll{{.Model.Name}}s() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			skip := c.Param("skip")
			limit := c.Param("limit")

			if skip == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "skip parameter has not been provided",
				})
				c.Abort()
				return
			}

			if limit == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "limit parameter has not been provided",
				})
				c.Abort()
				return
			}

			iskip, err := strconv.ParseInt(skip, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err,
				})
				c.Abort()
				return
			}

			ilimit, err := strconv.ParseInt(limit, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err,
				})
				c.Abort()
				return
			}
			selector := make(map[string]interface{})
			jsonMap := c.Request.URL.Query()

			for key, val := range jsonMap {
				selector[key] = val[0]
			}

			{{.Model.Name | ToLower}}s, err := {{.Model.Name | Initial}}.I{{.Model.Name}}.GetAll{{.Model.Name}}s(int64(iskip), int64(ilimit), selector)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, {{.Model.Name | ToLower}}s)
		}
	}
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}) GetAll{{.Model.Name}}sBy() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			skip := c.Param("skip")
			limit := c.Param("limit")

			if skip == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "skip parameter has not been provided",
				})
				c.Abort()
				return
			}

			if limit == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "limit parameter has not been provided",
				})
				c.Abort()
				return
			}

			iskip, err := strconv.ParseInt(skip, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err,
				})
				c.Abort()
				return
			}

			ilimit, err := strconv.ParseInt(limit, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err,
				})
				c.Abort()
				return
			}

			qstring := c.Request.URL.Query().Get("search")

			selector := make(map[string]interface{})
			jsonMap := c.Request.URL.Query()

			for key, val := range jsonMap {
				if key != "search" {
					selector[key] = val[0]
				}
			}
			{{.Model.Name | ToLower}}s, err := {{.Model.Name | Initial}}.I{{.Model.Name}}.GetAll{{.Model.Name}}sBy(qstring, selector, int64(iskip), int64(ilimit))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			//c.BindJSON(&profiles)
			c.JSON(http.StatusOK, {{.Model.Name | ToLower}}s)
		}
	}
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}) Update{{.Model.Name}}ByID() func(c *gin.Context) {
	var err error
	return func(c *gin.Context) {
		if c.Request.Method == "PUT" {

			id := c.Param("id")
			if id == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "id parameter has not been provided",
				})
				c.Abort()
				return
			}

			var {{.Model.Name | ToLower}} map[string]interface{}
			{{.Model.Name | ToLower}} = make(map[string]interface{})

			err = json.NewDecoder(c.Request.Body).Decode(&{{.Model.Name | ToLower}})
			
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			result, err := {{.Model.Name | Initial}}.I{{.Model.Name}}.Update{{.Model.Name}}ByID(id, {{.Model.Name | ToLower}})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": string(result.(int64)),
			})
			c.Abort()
			return
		}
	}
}

func ({{.Model.Name | Initial}} *{{.Model.Name}}) Delete{{.Model.Name}}ByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.Method == "DELETE" {

			id := c.Param("id")
			if id == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": "id parameter has not been provided",
				})
				c.Abort()
				return
			}

			result, err := {{.Model.Name | Initial}}.I{{.Model.Name}}.Delete{{.Model.Name}}ByID(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": string(result.(int64)),
			})
			c.Abort()
			return
		}
	}
}
`
	tm["interfaces"] = `// Author : readyGo "JitenP@Outlook.Com"
// This code is generated by template engine. You are free to make amendments as and where required
package interfaces

import (
	"{{.Project}}/models"
)


// {{.Model.Name}} Interface is to define all {{.Model.Name | ToLower}} related methods
type {{.Model.Name}}Interface interface {
    Create{{.Model.Name}}({{.Model.Name | ToLower}} *models.{{.Model.Name}})(interface{}, error)
	Update{{.Model.Name}}ByID(id string, data map[string]interface{}) (result interface{}, err error)
	Delete{{.Model.Name}}ByID(id string) (result interface{}, err error)
	Get{{.Model.Name}}ByID(id string) (*models.{{.Model.Name}}, error)
	GetAll{{.Model.Name}}s(skip int64, limit int64, selector interface{}) ([]models.{{.Model.Name}}, error)
	GetAll{{.Model.Name}}sBy(search string, selector interface{}, skip int64, limit int64) ([]models.{{.Model.Name}}, error)
    // Write any additional methods and implement them in database package
}
`
	tm["models"] = `// Author : readyGo "JitenP@Outlook.Com"
// This code is generated by template engine. You are free to make amendments as and where required
package models
{{- $validation := ""}}
{{- range .Model.Fields}}
{{- if .ValidateExp }}
{{- $validation = "true"}}
{{- end}}
{{- end}}
import (
	{{- if eq $validation "true" }}
	"errors"
	"regexp"
	{{- end}}
)
	type {{ .Model.Name }} struct {
	{{- range .Model.Fields }}
	{{ .Name }} {{ .Type }}` + "`" + `json:"{{.Name | ToLower}}"  {{ if eq $.config.DB "mongo" }}{{ if eq .Name "Id"}}bson:"_id,omitempty"{{else}}bson:"{{.Name | ToLower}}"{{- end}}{{- end}}` + "`" + `
	{{- end}}
	}
{{ $ModelName := .Model.Name }}
func Validate{{.Model.Name}}({{.Model.Name | Initial}} *{{.Model.Name}})(err error){
{{- range .Model.Fields}}
{{- if .ValidateExp }}
var rx{{.Name}}=regexp.MustCompile("{{.ValidateExp}}")
if !rx{{.Name}}.MatchString({{$ModelName | Initial}}.{{.Name}}) {
	return errors.New({{$ModelName | Initial}}.{{.Name}} + " is not a valid value for given field {{.Name}}")
}	
{{- end}}
{{- end}}
	return nil
}
`
	tm["main"] = `// Author : readyGo "JitenP@Outlook.Com"
// This code is generated by template engine. You are free to make amendments as and where required
package main

import (
	"{{$.config.Project}}/database"
	"{{$.config.Project}}/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	DBConnection = "mongodb://localhost:27017"
	DBName       = "demoDb"
)

func main() {
log.Println("Application {{.Project}} has started")

session, err := database.GetConnection(DBConnection, DBName)

	if err != nil {
		log.Fatal("mongodb database is not connected", err)
	}
	log.Println(session)

	gin.ForceConsoleColor()

	router := gin.Default()

     {{ range $.config.Models }}
     {{.Name | ToLower}} := new(handlers.{{.Name}})
	 {{.Name | ToLower}}.I{{.Name}} = &database.{{.Name}}DB{DB: session}
     {{.Name | ToLower}}Group := router.Group("/v1/{{.Name | ToLower}}")
    {
     {{.Name | ToLower}}Group.POST("/create", {{.Name | ToLower}}.Create{{.Name}}())
	 {{.Name | ToLower}}Group.DELETE("/delete/:id", {{.Name | ToLower}}.Delete{{.Name}}ByID())
	 {{.Name | ToLower}}Group.PUT("/update/:id", {{.Name | ToLower}}.Update{{.Name}}ByID())
	 {{.Name | ToLower}}Group.GET("/get/:id", {{.Name | ToLower}}.Get{{.Name}}ByID())
	 {{.Name | ToLower}}Group.GET("/getAll/:skip/:limit", {{.Name | ToLower}}.GetAll{{.Name}}s())
	 {{.Name | ToLower}}Group.GET("/getAllBy/:skip/:limit", {{.Name | ToLower}}.GetAll{{.Name}}sBy())
    }
     {{end}}
    router.Run(":{{$.config.Port}}")

}
`
	return tm
}
