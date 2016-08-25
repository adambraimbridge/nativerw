package main

import (
	"encoding/json"
	"io"
	"os"
)

type server struct {
	Port string `json:"port"`
}

type configuration struct {
	Mongos      string   `json:"mongos"`
	DbName      string   `json:"dbName"`
	Server      server   `json:"server"`
	Collections []string `json:"collections"`
}

func readConfigFromReader(r io.Reader) (c *configuration, e error) {
	c = new(configuration)

	decoder := json.NewDecoder(r)
	e = decoder.Decode(c)
	if e != nil {
		return nil, e
	}

	return
}

func readConfig(confPath string) (c *configuration, e error) {
	file, fErr := os.Open(confPath)
	defer file.Close()
	if fErr != nil {
		return nil, fErr
	}
	return readConfigFromReader(file)
}
