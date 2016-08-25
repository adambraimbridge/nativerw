package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
)

const txHeaderKey = "X-Request-Id"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const txHeaderLength = 20

var uuidRegexp = regexp.MustCompile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")

func (ma *mgoAPI) readContent(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	resourceID := vars["resource"]
	collection := vars["collection"]
	ctxlogger := txCombinedLogger{logger, obtainTxID(req)}

	if err := ma.validateAccess(collection, resourceID); err != nil {
		msg := fmt.Sprintf("Invalid collectionId (%v) or resourceId (%v).\n%v", collection, resourceID, err)
		ctxlogger.info(msg)
		http.Error(writer, msg, http.StatusBadRequest)
		return
	}

	found, resource, err := ma.Read(collection, resourceID)
	if err != nil {
		msg := fmt.Sprintf("Reading from mongoDB failed.\n%v\n", err.Error())
		ctxlogger.error(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	}
	if !found {
		msg := fmt.Sprintf("Resource not found. collection: %v, id: %v", collection, resourceID)
		ctxlogger.warn(msg)

		writer.Header().Add("Content-Type", "application/json")
		respBody, _ := json.Marshal(map[string]string{"message": msg})
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprint(writer, string(respBody))
		return
	}

	writer.Header().Add("Content-Type", resource.ContentType)

	om := outMappers[resource.ContentType]
	if om == nil {
		msg := fmt.Sprintf("Unable to handle resource of type %T. resourceId: %v, resource: %v", resource, resourceID, resource)
		ctxlogger.warn(msg)
		http.Error(writer, msg, http.StatusNotImplemented)
		return
	}
	err = om(writer, resource)
	if err != nil {
		msg := fmt.Sprintf("Unable to extract native content from resource with id %v. %v", resourceID, err.Error())
		ctxlogger.warn(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
	} else {
		ctxlogger.info(fmt.Sprintf("Read native content. resource_id: %+v", resourceID))
	}
}

func (ma *mgoAPI) getIds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coll := vars["collection"]
	ctxLogger := txCombinedLogger{logger, obtainTxID(r)}

	if err := ma.validateAccessForCollection(coll); err != nil {
		msg := fmt.Sprintf("Invalid collectionId (%v).\n%v", coll, err)
		ctxLogger.info(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)
	stop := make(chan struct{})
	defer close(stop)
	all, err := ma.Ids(coll, stop)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := struct {
		ID string `json:"id"`
	}{}
	for docID := range all {
		id.ID = docID
		enc.Encode(id)
	}
}

type outMapper func(io.Writer, resource) error

var outMappers = map[string]outMapper{
	"application/json": func(w io.Writer, resource resource) error {
		encoder := json.NewEncoder(w)
		return encoder.Encode(resource.Content)
	},
	"application/octet-stream": func(w io.Writer, resource resource) error {
		data := resource.Content.([]byte)
		_, err := io.Copy(w, bytes.NewReader(data))
		return err
	},
}

func (ma *mgoAPI) writeContent(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	collectionID := mux.Vars(req)["collection"]
	resourceID := mux.Vars(req)["resource"]
	ctxlogger := txCombinedLogger{logger, obtainTxID(req)}

	if err := ma.validateAccess(collectionID, resourceID); err != nil {
		msg := fmt.Sprintf("Invalid collectionId (%v) or resourceId (%v).\n%v", collectionID, resourceID, err)
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

	wrappedContent := wrap(content, resourceID, contentType)

	if err := ma.Write(collectionID, wrappedContent); err != nil {
		msg := fmt.Sprintf("Writing to mongoDB failed:\n%v\n", err)
		ctxlogger.error(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	}

	ctxlogger.info(fmt.Sprintf("Written native content. resource_id: %+v", resourceID))
}

func (ma *mgoAPI) validateAccess(collectionID, resourceID string) error {
	if ma.collections[collectionID] && uuidRegexp.MatchString(resourceID) {
		return nil
	}
	return errors.New("Collection not supported or resourceId not a valid uuid.")
}

func (ma *mgoAPI) validateAccessForCollection(collectionID string) error {
	if ma.collections[collectionID] {
		return nil
	}
	return errors.New("Collection not supported.")
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

func wrap(content interface{}, resourceID, contentType string) resource {
	return resource{
		UUID:        resourceID,
		Content:     content,
		ContentType: contentType,
	}
}

func obtainTxID(req *http.Request) string {
	txID := req.Header.Get(txHeaderKey)
	if txID == "" {
		return randSeq(txHeaderLength)
	}
	return txID
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func testSlowRequest(writer http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Millisecond * 1100)
}
