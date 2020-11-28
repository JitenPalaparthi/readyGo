package mapping

import (
	"encoding/json"
	"errors"
)

// Reader interface reads a file and retunrs it as string
type Reader interface {
	Read(string) (string, error)
}

// Mapping contains mapping information
type Mapping struct {
	Key     string
	OpsData []OpsData
	Reader  Reader
}

// OpsData contains what is the source and destination
type OpsData struct {
	OpType string `json:"opType"`
	Src    string `json:"src`
	Dst    string `json:"dst`
}

// New creates a new Mapping
func New(reader Reader, file, key string) (mapping *Mapping, err error) {
	if reader == nil {
		return nil, errors.New("reader object is not provided or nil")
	}
	content, err := reader.Read(file) // Read the mappings file
	if err != nil {
		return nil, err
	}
	fullOpsMap := make(map[string][]OpsData)
	err = json.Unmarshal([]byte(content), &fullOpsMap)
	if err != nil {
		return nil, err
	}
	opsData, ok := fullOpsMap[key]
	if !ok {
		return nil, errors.New("required data not found")
	}
	mapping = &Mapping{}
	mapping.Key = key
	mapping.OpsData = opsData
	mapping.Reader = reader
	return mapping, err
}
