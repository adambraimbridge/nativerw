package main

import (
    "encoding/json"
    "os"
    "io"
)

type Server struct {
	Port	string	`json: port`
}

type Configuration struct {
	DbName	string `json: dbName`
	Server  Server `json: server`
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
	file, fErr := os.Open("conf.json")
	if (fErr != nil) {
		return nil, fErr
	}
	return readConfigFromReader(file)
}