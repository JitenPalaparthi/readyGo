package mapping

import (
	"encoding/json"
	"errors"
)

var (
	// ErrReaderNotProvided is to return when reader is nil or empty
	ErrReaderNotProvided = errors.New("reader is empty or nil")
	// ErrNoData is to return err when there is not data
	ErrNoData = errors.New("no data found")
)

// Reader interface reads a file and retunrns it as string
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
	Key        string `json:"key"`
	OpType     string `json:"opType"`
	Src        string `json:"src"`
	Dst        string `json:"dst"`
	Ext        string `json:"ext"`
	GenForType string `json:"genForType"`
}

// New creates a new Mapping
func New(reader Reader, file, key string) (mapping *Mapping, err error) {
	if reader == nil {
		return nil, ErrReaderNotProvided
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
	opsData, ok := fullOpsMap["data"]
	if !ok {
		return nil, ErrNoData
	}
	mapping = &Mapping{}
	mapping.Key = key
	mapping.OpsData = opsData
	mapping.Reader = reader
	return mapping, err
}
