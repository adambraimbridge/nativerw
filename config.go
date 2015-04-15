package main

import (
    "encoding/json"
    "os"
)

type Server struct {
	Port	string	`json: port`
}

type Configuration struct {
	DbName	string `json: dbName`
	Server  Server `json: server`
}

func readConfig() (c *Configuration, e error) {
	file, fErr := os.Open("conf.json")
	if (fErr != nil) {
		return nil, fErr
	}

	c = new(Configuration)

	decoder := json.NewDecoder(file)
	e = decoder.Decode(c)
	if e != nil {
		return nil, e
	}

	return 
}