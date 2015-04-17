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

	unwrappedResource := unwrapResource(resource)
	contentType := getContentType(resource)

	writer.Header().Add("Content-Type", contentType)
	encoder := json.NewEncoder(writer)
	encoder.Encode(unwrappedResource)
}

func getContentType(resource interface{}) string {
	return resource.(map[string]interface{})["content-type"].(string)
}

func unwrapResource(resource interface{}) interface{} {
	return resource.(map[string]interface{})["content"]
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

func extractContent(req *http.Request, resourceId string) (map[string]interface{}, error) {
	var wrappedContent map[string]interface{}
	var err error
	if req.Header.Get("Content-Type") == "application/json" {
		var content map[string]interface{}
		content, err = extractJson(req, resourceId)
		wrappedContent = wrapMap(content, resourceId, "application/json")
	} else {
		var binary []byte
		binary, err = extractBinary(req)
		wrappedContent = wrapBinary(binary, resourceId, "application/octet-stream")
	}
	if err != nil {
		return nil, err
	}
	return wrappedContent, nil
}

func extractJson(req *http.Request, resourceId string) (map[string]interface{}, error) {
	var content map[string]interface{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&content); err != nil {
		return nil, &extractionError{fmt.Sprintf("JSON decode failed:\n%v\n", err), http.StatusBadRequest}
	}
	if payloadId := content[uuidName]; payloadId != resourceId {
		return nil, &extractionError{fmt.Sprintf("Given resource id %v does not match id in payload %v .",
			resourceId, payloadId), http.StatusBadRequest}
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

func wrapMap(content map[string]interface{}, resourceId, contentType string) map[string]interface{} {
	return map[string]interface{}{
		"uuid":         resourceId,
		"content":      content,
		"content-type": contentType,
	}
}

func wrapBinary(content []byte, resourceId, contentType string) map[string]interface{} {
	return map[string]interface{}{
		"uuid":         resourceId,
		"content":      content,
		"content-type": contentType,
	}
}
