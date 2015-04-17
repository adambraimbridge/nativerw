package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

func (ma *MgoApi) readContent(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	resourceId := vars["resource"]
	collection := vars["collection"]

	found, resource := ma.Read(collection, resourceId)
	if !found {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf("resource with id %s was not found\n", resourceId)))
		return
	}

	writer.Header().Add("Content-Type", resource.ContentType)
	encoder := json.NewEncoder(writer)
	encoder.Encode(resource.Content)
}

func (mgoApi *MgoApi) writeContent(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	collectionId := mux.Vars(req)["collection"]
	resourceId := mux.Vars(req)["resource"]

	content, contentType, err := extractContent(req)
	if err != nil {
		// TODO: this could be a server error too?
		http.Error(writer, fmt.Sprintf("Extracting content from HTTP body failed:\n%v\n", err), http.StatusBadRequest)
		return
	}

	wrappedContent := wrap(content, resourceId, contentType)

	if err := mgoApi.Write(collectionId, wrappedContent); err != nil {
		http.Error(writer, fmt.Sprintf("Writing to mongoDB failed:\n%v\n", err), http.StatusInternalServerError)
		return
	}
}

func extractContent(req *http.Request) (content interface{}, contentType string, err error) {
	contentType = req.Header.Get("Content-Type")
	mapper := inMappers[contentType]
	if mapper == nil {
		// default to binary
		contentType = "application/octet-stream"
		mapper = inMappers[contentType]
	}
	content, err = mapper(req.Body)
	return
}

type inMapper func(io.Reader) (interface{}, error)

var inMappers = map[string]inMapper{
	"application/json": func(r io.Reader) (interface{}, error) {
		var c map[string]interface{}
		err := json.NewDecoder(r).Decode(&c)
		return c, err
	},
	"application/octet-stream": func(r io.Reader) (interface{}, error) {
		return ioutil.ReadAll(r)
	},
}

func wrap(content interface{}, resourceId, contentType string) Resource {
	return Resource{
		UUID:        resourceId,
		Content:     content,
		ContentType: contentType,
	}
}
