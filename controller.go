package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

type extractionError struct {
	cause    string
	httpCode int
}

func (e extractionError) Error() string {
	return e.cause
}

func (mgoApi *MgoApi) writeContent(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	collectionId := mux.Vars(req)["collection"]
	resourceId := mux.Vars(req)["resource"]

	wrappedContent, exErr := extractContent(req, resourceId)
	if exErr != nil {
		err := exErr.(*extractionError)
		http.Error(writer, fmt.Sprintf("Extracting content from HTTP body failed:\n%v\n", exErr), err.httpCode)
		return
	}

	if wrErr := mgoApi.Write(collectionId, wrappedContent); wrErr != nil {
		http.Error(writer, fmt.Sprintf("Writing to mongoDB failed:\n%v\n", wrErr), http.StatusInternalServerError)
		return
	}
}

func extractContent(req *http.Request, resourceId string) (Resource, error) {
	var err error
	var content interface{}
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/json" {
		content, err = extractJson(req)
	} else {
		content, err = extractBinary(req)
		contentType = "application/octet-stream"
	}
	if err != nil {
		return Resource{}, err
	}
	return wrap(content, resourceId, contentType), nil
}

func extractJson(req *http.Request) (map[string]interface{}, error) {
	var content map[string]interface{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&content); err != nil {
		return nil, &extractionError{fmt.Sprintf("JSON decode failed:\n%v\n", err), http.StatusBadRequest}
	}
	return content, nil
}

func extractBinary(req *http.Request) ([]byte, error) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []byte{}, &extractionError{fmt.Sprintf("Reading the body of the request failed:\n%v\n", err),
			http.StatusInternalServerError}
	}
	return content, nil
}

func wrap(content interface{}, resourceId, contentType string) Resource {
	return Resource{
		UUID:        resourceId,
		Content:     content,
		ContentType: contentType,
	}
}
