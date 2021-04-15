package generate

import (
	"readyGo/mapping"
	"readyGo/scaler"
)

// Implementer is an interface used to implement validation and lang specific logic
type Implementer interface {
	IsValidIdentifier(fielden string) bool
	///SetFieldCategory()
	//IsModelType(iden string) bool
	GetFuncReturnType(interface{}) string
}

// Generate is a type that holds configuration data
type Generate struct {
	Version       string        `json:"version" yaml:"version"`
	Project       string        `json:"project" yaml:"project"`             // ideally project root directory .i.e project name
	Kind          string        `json:"kind" yaml:"kind"`                   // Kind of the project http , grpc , CloudEvents , cli
	Port          string        `json:"port" yaml:"port"`                   // Port that is used to communicate http project
	APISpec       APISpec       `json:"apiSpec" yaml:"apiSpec"`             // Api releted information generally used to design apis
	DatabaseSpec  DatabaseSpec  `json:"databaseSpec" yaml:"databaseSpec"`   // Database related information like sql|mongo connection string and db name .. to be maintained in this
	MessagingSpec MessagingSpec `json:"messagingSpec" yaml:"messagingSpec"` // Messaging related information. kind is nsq | nats | kafka
	Models        []Model       `json:"models" yaml:"models"`
	Mapping       *mapping.Mapping
	Scalers       scaler.Map
	Implementer   Implementer // interface to use lang specific implementation logic
}

// Model is to hold model data from configuration file
type Model struct {
	Name               string             `json:"name" yaml:"name"`
	MessagingModelSpec MessagingModelSpec `json:"messagingModelSpec" yaml:"messagingModelSpec"`
	Fields             []Field            `json:"fields" yaml:"fields"`
}

// MessagingModelSpec is to define model specific messging metadata
type MessagingModelSpec struct {
	MessageRespondType string `json:"messageRespondType" yaml:"messageRespondType"` // There are two types as of now. publish|subscribe
	Topic              string `json:"topic" yaml:"topic"`                           // Topic is topic for nats , subject for kafka
}

// Field is to hold fields in a model that comes from configuration file
type Field struct {
	Name        string `json:"name" yaml:"name"`               // Name of the field. Should be valid Go identifier
	Type        string `json:"type" yaml:"type"`               // type should be either readyGo specified scaler type or a defined model
	Definition  string `json:"definition" yaml:"definition"`   // definition is to define function.Dont want to give more fields to user to set so type should be auto matially taken when type is given as function
	IsKey       bool   `json:"isKey" yaml:"isKey"`             // If it is a key field . Key fields generates different methods to check the data in the database is unique or not
	ValidateExp string `json:"validateExp" yaml:"validateExp"` // Regular expression that would be used for field level validations in models
	Category    string `json:"category" yaml:"category"`       // Category is general field types are scalers if its a type for a model then it is model etc
	Annotation  string `json:"annotation" yaml:"annotation"`   // Field Annotations are generally used for json bson gorm etc.
}

// DatabaseSpec struct type contains database related information
type DatabaseSpec struct {
	Kind             string `json:"kind" yaml:"kind"` //sql nosql
	Name             string `json:"name" yaml:"name"` // mongo mysql
	ConnectionString string `json:"connectionString" yaml:"connectionString"`
	DBName           string `json:"dbName" yaml:"dbName"`
}

// APISpec struct type contains api related information
type APISpec struct {
	Kind    string `json:"kind" yaml:"kind"`       // http | grpc | cloudEvent
	Port    string `json:"port" yaml:"port"`       // port to run on
	Version string `json:"version" yaml:"version"` // Version that is used to define apis.example v1/public/get v2/private/create etc.
}

// MessagingSpec struct type contains message queue related information
type MessagingSpec struct {
	Kind             string `json:"kind" yaml:"kind"`
	ConnectionString string `json:"connectionString" yaml:"connectionString"`
}
