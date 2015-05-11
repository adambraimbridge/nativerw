package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"errors"
)

const txHeaderKey = "X-Request-Id"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const txHeaderLength = 20
var uuidRegexp = regexp.MustCompile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")

func (ma *MgoApi) readContent(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	resourceId := vars["resource"]
	collection := vars["collection"]
	ctxlogger := TxCombinedLogger{logger, obtainTxId(req)}

	if err := ma.validateAccess(collection, resourceId); err != nil {
		msg := fmt.Sprintf("Invalid collectionId (%v) or resourceId (%v).\n%v", collection, resourceId, err)
		ctxlogger.info(msg)
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	found, resource, err := ma.Read(collection, resourceId)
	if err != nil {
		msg := fmt.Sprintf("Reading from mongoDB failed.\n%v\n", err.Error())
        ctxlogger.error(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	}
	if !found {
		msg := fmt.Sprintf("Resource not found. collection: %v, id: %v", collection, resourceId)
		ctxlogger.warn(msg)
		http.Error(writer, msg, http.StatusNotFound)
		return
	}

	writer.Header().Add("Content-Type", resource.ContentType)

	om := outMappers[resource.ContentType]
	if om == nil {
		msg := fmt.Sprintf("Unable to handle resource of type %T. resourceId: %v, resource: %v", resource, resourceId, resource)
		ctxlogger.warn(msg)
		http.Error(writer, msg, http.StatusNotImplemented)
		return
	}
	err = om(writer, resource)
	if err != nil {
		msg := fmt.Sprintf("Unable to extract native content from resource with id %v. %v", resourceId, err.Error())
		ctxlogger.warn(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
	} else {
		ctxlogger.info(fmt.Sprintf("Read native content. resource_id: %+v", resourceId))
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
	ctxlogger := TxCombinedLogger{logger, obtainTxId(req)}

	if err := mgoApi.validateAccess(collectionId, resourceId); err != nil {
		msg := fmt.Sprintf("Invalid collectionId (%v) or resourceId (%v).\n%v", collectionId, resourceId, err)
		ctxlogger.info(msg)
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	contentType := req.Header.Get("Content-Type")
	mapper := inMappers[contentType]
	if mapper == nil {
		msg := fmt.Sprintf("Content-Type header missing. Default value ('application/octet-stream') is used.")
		ctxlogger.info(msg)
		contentType = "application/octet-stream"
		mapper = inMappers[contentType]
	}

	content, err := mapper(req.Body)
	if err != nil {
		msg := fmt.Sprintf("Extracting content from HTTP body failed:\n%v\n", err)
		ctxlogger.warn(msg)
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	wrappedContent := wrap(content, resourceId, contentType)

	if err := mgoApi.Write(collectionId, wrappedContent); err != nil {
		msg := fmt.Sprintf("Writing to mongoDB failed:\n%v\n", err)
		ctxlogger.error(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	} else {
		ctxlogger.info(fmt.Sprintf("Written native content. resource_id: %+v", resourceId))
	}
}

func (mgoApi *MgoApi) validateAccess(collectionId, resourceId string) error {
	if mgoApi.collections[collectionId] && uuidRegexp.MatchString(resourceId) {
		return nil
	}
	return errors.New("Collection not supported or resourceId not a valid uuid.")
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

func obtainTxId(req *http.Request) string {
	txId := req.Header.Get(txHeaderKey)
	if txId == "" {
		return randSeq(txHeaderLength)
	}
	return txId
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
