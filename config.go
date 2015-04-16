package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type Mongo struct {
	Host string `json: host`
	Port string `json: port`
}

type Server struct {
	Port string `json: port`
}

type Configuration struct {
	Mongos []Mongo `json: mongos`
	DbName string  `json: dbName`
	Server Server  `json: server`
}

func (c *Configuration) prepareMgoUrls() string {
	var hostsPorts string
	for _, mongo := range c.Mongos {
		hostsPorts += mongo.Host + ":" + mongo.Port + ","
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

func readConfig() (c *Configuration, e error) {
	file, fErr := os.Open("config.json")
	if fErr != nil {
		return nil, fErr
	}
	return readConfigFromReader(file)
}
