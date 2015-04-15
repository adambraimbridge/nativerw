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
)

type MgoApi struct {
	dbName         string
	session        *mgo.Session
	resourceIdName string
	beforeWrite    propertyConverter
	afterRead      propertyConverter
}

func NewMgoApi(dbName string, resourceIdName string, beforeWrite, afterRead propertyConverter) (*MgoApi, error) {
	s, err := mgo.DialWithTimeout("localhost", time.Duration(3*time.Second))
	if err != nil {
		return nil, err
	}
	s.SetMode(mgo.Monotonic, true)

	return &MgoApi{dbName, s, resourceIdName, beforeWrite, afterRead}, nil
}

func (ma *MgoApi) Write(collection string, resource map[string]interface{}) error {
	coll := ma.session.DB(ma.dbName).C(collection)

	ma.mongoizeAll(resource)

	_, err := coll.Upsert(bson.D{{ma.resourceIdName, resource[ma.resourceIdName]}}, resource)

	return err
}

func (ma *MgoApi) Read(collection string, resourceId string) (bool, interface{}) {
	coll := ma.session.DB(ma.dbName).C(collection)

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
	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(resource)
}

func (ma *MgoApi) writeHandler(w http.ResponseWriter, req *http.Request) {
	collectionId := mux.Vars(req)["collection"]
	resourceId := mux.Vars(req)["resource"]

	dec := json.NewDecoder(req.Body)
	var resource map[string]interface{}

	if err := dec.Decode(&resource); err != nil {
		http.Error(w, fmt.Sprintf("json decode failed: %v", err.Error()), http.StatusBadRequest)
		return
	}

	if plResourceId := resource[ma.resourceIdName]; plResourceId != resourceId {
		http.Error(w, "given resource id does not match payload", http.StatusBadRequest)
		return 
	}

	wrappedResource := wrapResource(resource, resourceId, "application/json")

	if err := ma.Write(collectionId, wrappedResource); err != nil {
		http.Error(w, fmt.Sprintf("write failed:\n%v\n", err), http.StatusInternalServerError)
		return
	}
}

func wrapResource(resource map[string]interface{}, resourceId, contentType string) map[string]interface{} {
	return map[string]interface{}{
    "uuid": resourceId,
    "content": resource,
    "content-type": contentType,
  }
}

func createMgoApi() (*MgoApi, error) {
	mgoApi, err := NewMgoApi(config.DbName, "uuid",
		compositePropertyConverter{[]propertyConverter{UUIDToBson, DateToBson}}.convert,
		compositePropertyConverter{[]propertyConverter{UUIDFromBson, DateFromBson, MongoIdRemover}}.convert,
	)
	return mgoApi, err
}

var config, configErr = readConfig()
var mgoApi, mgoApiCreationErr = createMgoApi()

func main() {

	if ( configErr != nil) {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", configErr.Error())
		return
	}

	if mgoApiCreationErr != nil {
		fmt.Fprintf(os.Stderr, "%v\n", mgoApiCreationErr.Error())
		return
	}

	m := mux.NewRouter()
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, m))
	m.HandleFunc("/{collection}/{resource}", mgoApi.readHandler).Methods("GET")
	m.HandleFunc("/{collection}/{resource}", mgoApi.writeHandler).Methods("PUT")
	m.HandleFunc("/__health", fthealth.Handler("Dependent services healthceck",
	  "Checking connectivity and usability of dependent services: mongoDB and native-ingester.", mgoHealth))

	http.ListenAndServe(":" + config.Server.Port, nil)
}
