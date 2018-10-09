package config

import (
	"encoding/json"
	"io"
	"os"
)

// Server config struct
type Server struct {
	Port int `json:"port"`
}

// Configuration data
type Configuration struct {
	Mongos      string   `json:"mongos"`
	DbName      string   `json:"dbName"`
	Server      Server   `json:"server"`
	Collections []string `json:"collections"`
}

// ReadConfigFromReader reads config as a json stream from the given reader
func ReadConfigFromReader(r io.Reader) (c *Configuration, e error) {
	c = new(Configuration)

	decoder := json.NewDecoder(r)
	e = decoder.Decode(c)
	if e != nil {
		return nil, e
	}

	return c, nil
}

// ReadConfig reads config as a json file from the given path
func ReadConfig(confPath string) (c *Configuration, e error) {
	file, fErr := os.Open(confPath)
	defer file.Close()
	if fErr != nil {
		return nil, fErr
	}
	return ReadConfigFromReader(file)
}
