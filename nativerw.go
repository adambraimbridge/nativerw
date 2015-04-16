package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"reflect"
	"time"
	"git.svc.ft.com/scm/gl/fthealth.git"
	"strings"
	"io/ioutil"
)

const uuidName = "uuid"

type MgoApi struct {
	dbName         string
	session        *mgo.Session
	resourceIdName string
	beforeWrite    propertyConverter
	afterRead      propertyConverter
}

func NewMgoApi(urls, dbName, resourceIdName string, beforeWrite, afterRead propertyConverter) (*MgoApi, error) {
	s, err := mgo.DialWithTimeout(urls, time.Duration(3*time.Second))
	if err != nil {
		return nil, err
	}
	s.SetMode(mgo.Monotonic, true)

	return &MgoApi{dbName, s, resourceIdName, beforeWrite, afterRead}, nil
}

func prepareMgoUrls(mongos []Mongo) string {
	var hostsPorts string
	for _, mongo := range mongos {
		hostsPorts += mongo.Host + ":" + mongo.Port + ","
	}
	return strings.TrimRight(hostsPorts, ",")
}

func (ma *MgoApi) Write(collection string, resource map[string]interface{}) error {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	ma.mongoizeAll(resource)

	_, err := coll.Upsert(bson.D{{ma.resourceIdName, resource[ma.resourceIdName]}}, resource)

	return err
}

func (ma *MgoApi) Read(collection string, resourceId string) (bool, interface{}) {
	newSession := ma.session.Copy()
	defer newSession.Close()

	coll := newSession.DB(ma.dbName).C(collection)

	// convert resource id to mgo friendly form if needed
	props := make(map[string]interface{})
	props[ma.resourceIdName] = resourceId
	ma.mongoizeAll(props)
	mongoResourceId := props[ma.resourceIdName]

	var resource map[string]interface{}
	if err := coll.Find(bson.M{ma.resourceIdName: mongoResourceId}).One(&resource); err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		panic(err)
	}

	ma.unmongoizeAll(resource)

	return true, resource
}

func (ma *MgoApi) mongoizeAll(resource map[string]interface{}) {
	for k, v := range resource {
		if reflect.ValueOf(v).Type() == mapStrIfType {
			ma.mongoizeAll(v.(map[string]interface{}))
		} else {
			pm := simplePropertyModifier{resource, k}
			ma.beforeWrite(pm, resource, k, v)
		}
	}
}

var mapStrIfType = reflect.ValueOf(make(map[string]interface{})).Type()

func (ma *MgoApi) unmongoizeAll(resource map[string]interface{}) {
	for k, v := range resource {
		if reflect.ValueOf(v).Type() == mapStrIfType {
			ma.unmongoizeAll(v.(map[string]interface{}))
		} else {
			pm := simplePropertyModifier{resource, k}
			ma.afterRead(pm, resource, k, v)
		}
	}
}

func (ma *MgoApi) readHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceId := vars["resource"]
	collection := vars["collection"]

	found, resource := ma.Read(collection, resourceId)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("resource with id %s was not found\n", resourceId)))
		return
	}

	unwrappedResource := unwrapResource(resource)
	contentType := getContentType(resource)

	w.Header().Add("Content-Type", contentType)
	enc := json.NewEncoder(w)
	enc.Encode(unwrappedResource)
}

func getContentType(resource interface{}) string {
	return resource.(map[string]interface{})["content-type"].(string)
}

func unwrapResource(resource interface{}) interface{} {
	return resource.(map[string]interface{})["content"]
}

type extractionError struct {
	cause string
	httpCode int
}

func (e extractionError) Error() string {
	return e.cause
}

func (mgoApi *MgoApi) writeHandler(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	collectionId := mux.Vars(req)["collection"]
	resourceId := mux.Vars(req)["resource"]

	wrappedContent, exErr := extractContent(req, resourceId)
	if exErr != nil {
		err := exErr.(extractionError)
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
		content, err = extractJson(req, resourceId);
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
		return nil, &extractionError{ fmt.Sprintf("JSON decode failed:\n%v\n", err), http.StatusBadRequest }
	}
	if payloadId := content[mgoApi.resourceIdName]; payloadId != resourceId {
		return nil, &extractionError{ fmt.Sprintf("Given resource id %v does not match id in payload %v .",
			resourceId, payloadId), http.StatusBadRequest }
	}
	return content, nil
}

func extractBinary(req *http.Request) ([]byte, error) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []byte {}, &extractionError{ fmt.Sprintf("Reading the body of the request failed:\n%v\n", err),
			http.StatusInternalServerError }
	}
	return content, nil
}

func wrapMap(content map[string]interface{}, resourceId, contentType string) map[string]interface{} {
	return map[string]interface{}{
		"uuid": resourceId,
		"content": content,
		"content-type": contentType,
	}
}

func wrapBinary(content []byte, resourceId, contentType string) map[string]interface{} {
	return map[string]interface{}{
		"uuid": resourceId,
		"content": content,
		"content-type": contentType,
	}
}

func createMgoApi(config *Configuration) (*MgoApi, error) {
	mgoApi, err := NewMgoApi(mongoUrls, config.DbName, uuidName,
		compositePropertyConverter{[]propertyConverter{UUIDToBson, DateToBson}}.convert,
		compositePropertyConverter{[]propertyConverter{UUIDFromBson, DateFromBson, MongoIdRemover}}.convert,
	)
	return mgoApi, err
}

var config, configErr = readConfig()
var mongoUrls = prepareMgoUrls(config.Mongos)
var mgoApi, mgoApiCreationErr = createMgoApi(config)

func main() {

	if configErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", configErr.Error())
		return
	}

	if mgoApiCreationErr != nil {
		fmt.Fprintf(os.Stderr, "%v\n", mgoApiCreationErr.Error())
		return
	}

	router := mux.NewRouter()
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, router))
	router.HandleFunc("/{collection}/{resource}", mgoApi.readHandler).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", mgoApi.writeHandler).Methods("PUT")
	router.HandleFunc("/__health", fthealth.Handler("Dependent services healthceck",
	  "Checking connectivity and usability of dependent services: mongoDB and native-ingester.", mgoHealth))

	http.ListenAndServe(":" + config.Server.Port, nil)
}
