package models

import (
	"bytes"
	"encoding/gob"
)

// ToBytes is to convert address to bytes
func ToBytes(o interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(o)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
