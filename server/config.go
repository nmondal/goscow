package server

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type VerbType string

const (
	CONNECT = "connect"
	DELETE  = "delete"
	GET     = "get"
	HEAD    = "head"
	OPTIONS = "options"
	PATCH   = "patch"
	POST    = "post"
	PUT     = "put"
	TRACE   = "trace"
)

const (
	JSON        = "json"
	JSON_PRETTY = "jsonp"
	XML         = "xml"
	XML_PRETTY  = "xmlp"
)

type GoSCowConfig struct {
	Base   string
	Name   string
	Port   int16
	Static string
	Reload bool
	Routes map[VerbType]map[string]string
}

func From(configFile string) (*GoSCowConfig, error) {
	path, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Dir(path) + "/"
	name := filepath.Base(path)
	data, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		return nil, err
	}
	configString := string(data)
	configString = strings.ReplaceAll(configString, "_/", baseDir)
	goSCowConfig := GoSCowConfig{}
	err = yaml.Unmarshal([]byte(configString), &goSCowConfig)
	if err != nil {
		return nil, err
	}
	goSCowConfig.Base = baseDir
	goSCowConfig.Name = name
	return &goSCowConfig, nil
}
