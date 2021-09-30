package scaler

import (
	"encoding/json"
	"errors"
)

var (
	// ErrReaderNotProvided is to return when reader is nil or empty
	ErrReaderNotProvided = errors.New("reader is empty or nil")
)

// Reader interface reads a file and retunrs it as string
type Reader interface {
	Read(string) (string, error)
}

// Scaler is to define types in the system
type Scaler struct {
	GoType   string
	GrpcType string
}

// Map is to hold key and value as &Scaler types
type Map map[string]*Scaler

// New creates a new ScalerType
func New(reader Reader, file string) (Map, error) {
	if reader == nil {
		return nil, ErrReaderNotProvided
	}
	content, err := reader.Read(file) // Read the scaler config file
	if err != nil {
		return nil, err
	}
	scalerType := make(map[string]*Scaler)
	err = json.Unmarshal([]byte(content), &scalerType)
	if err != nil {
		return nil, err
	}
	return Map(scalerType), err
}

// GetScaler is to fetch Scaler type from the map
func (m Map) GetScaler(configType string) *Scaler {
	value, ok := m[configType]
	if ok {
		return value
	}
	return nil
}

// IsValidreadyGotype is to check whether type is valid type or not
func (m Map) IsValidreadyGotype(configType string) bool {
	_, ok := m[configType]
	return ok
}
