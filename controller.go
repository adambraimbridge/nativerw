package main

import (
	"bytes"
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

	found, resource, err := ma.Read(collection, resourceId)
	if err != nil {
        msg := fmt.Sprintf("Reading from mongoDB failed.\n%v\n", err.Error())
        logger.warn(msg)
        http.Error(writer, msg, http.StatusInternalServerError)
        return
	}
	if !found {
        msg := fmt.Sprintf("Resource not found. collection: %v, id: %v\n", collection, resourceId)
        logger.warn(msg)
        http.Error(writer, msg, http.StatusNotFound)
		return
	}

	writer.Header().Add("Content-Type", resource.ContentType)

	om := outMappers[resource.ContentType]
	if om == nil {
        msg := fmt.Sprintf("Unable to handle resource of type %T. resourceId: %v, resource: %v\n", resource, resourceId, resource)
        logger.warn(msg)
        http.Error(writer, msg, http.StatusNotImplemented)
        return
	}
	err = om(writer, resource)
	if err != nil {
        msg := fmt.Sprintf("Unable to extract native content from resource with id %v. %v\n", resourceId, err.Error())
        logger.warn(msg)
        http.Error(writer, msg, http.StatusInternalServerError)
	}
}

type outMapper func(io.Writer, Resource) error

var outMappers = map[string]outMapper{
	"application/json": func(w io.Writer, resource Resource) error {
		encoder := json.NewEncoder(w)
		return encoder.Encode(resource.Content)
	},
	"application/octet-stream": func(w io.Writer, resource Resource) error {
		data := resource.Content.([]byte)
		_, err := io.Copy(w, bytes.NewReader(data))
		return err
	},
}

func (mgoApi *MgoApi) writeContent(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	collectionId := mux.Vars(req)["collection"]
	resourceId := mux.Vars(req)["resource"]

	contentType := req.Header.Get("Content-Type")
	mapper := inMappers[contentType]
	if mapper == nil {
		// default to binary
		contentType = "application/octet-stream"
		mapper = inMappers[contentType]
	}

	content, err := mapper(req.Body)
	if err != nil {
		// TODO: this could be a server error too?
        msg := fmt.Sprintf("Extracting content from HTTP body failed:\n%v\n", err)
        logger.warn(msg)
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	wrappedContent := wrap(content, resourceId, contentType)

	if err := mgoApi.Write(collectionId, wrappedContent); err != nil {
        msg := fmt.Sprintf("Writing to mongoDB failed:\n%v\n", err)
        logger.warn(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	}
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
