package configure

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Configure is to configure file copy sets
type Configure struct {
	FC map[string][]StaticFileSet
	TC map[string][]string
	DC map[string][]string // Directory configurations
}

// StaticFileSet is to get static file sets
type StaticFileSet struct {
	Src, Dst string
}

// New to create new configure object
func New(file *string) (c *Configure, err error) {
	if file == nil || *file == "" {
		return nil, errors.New("no file provided")
	}
	ext := filepath.Ext(*file)
	fmt.Println(ext)
	if ext != ".json" {
		return nil, errors.New("Only json files are allowed ")
	}
	cFile, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cFile), &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ReadFC is to read file configurations
func (c *Configure) ReadFC(key string) []StaticFileSet {
	v, ok := c.FC[key]
	if ok {
		return v
	}
	return nil
}

// ReadTC is to read Template Configrations
func (c *Configure) ReadTC(key string) []string {
	v, ok := c.TC[key]
	if ok {
		return v
	}
	return nil
}

// ReadDC is to read directory Configrations
func (c *Configure) ReadDC(key string) []string {
	v, ok := c.DC[key]
	if ok {
		return v
	}
	return nil
}
