package scalar

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

// Scalar is to define types in the system
type Scalar struct {
	GoType   string
	GrpcType string
}

// Map is to hold key and value as &Scalar types
type Map map[string]*Scalar

// New creates a new ScalerType
func New(reader Reader, file string) (Map, error) {
	if reader == nil {
		return nil, ErrReaderNotProvided
	}
	content, err := reader.Read(file) // Read the scalar config file
	if err != nil {
		return nil, err
	}
	scalerType := make(map[string]*Scalar)
	err = json.Unmarshal([]byte(content), &scalerType)
	if err != nil {
		return nil, err
	}
	return Map(scalerType), err
}

// GetScalar is to fetch Scalar type from the map
func (m Map) GetScalar(configType string) *Scalar {
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
