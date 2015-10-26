package main

import (
	"encoding/json"
	"io"
	"os"
)

type Server struct {
	Port string `json:"port"`
}

type Configuration struct {
	Mongos      string   `json:"mongos"`
	DbName      string   `json:"dbName"`
	Server      Server   `json:"server"`
	Collections []string `json:"collections"`
}

func readConfigFromReader(r io.Reader) (c *Configuration, e error) {
	c = new(Configuration)

	decoder := json.NewDecoder(r)
	e = decoder.Decode(c)
	if e != nil {
		return nil, e
	}

	return
}

func readConfig(confPath string) (c *Configuration, e error) {
	file, fErr := os.Open(confPath)
	defer file.Close()
	if fErr != nil {
		return nil, fErr
	}
	return readConfigFromReader(file)
}
