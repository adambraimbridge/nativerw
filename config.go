package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type Server struct {
	Port string `json: port`
}

type Configuration struct {
	Mongos      []string `json: mongos`
	DbName      string   `json: dbName`
	Server      Server   `json: server`
	Collections []string `json: collections`
}

func (c *Configuration) prepareMgoUrls() string {
	var hostsPorts string
	for _, mongo := range c.Mongos {
		hostsPorts += mongo + ","
	}
	return strings.TrimRight(hostsPorts, ",")
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
