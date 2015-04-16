package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"git.svc.ft.com/scm/gl/fthealth.git"
)

func createMgoApi(config *Configuration) (*MgoApi, error) {
	mgoUrls := config.prepareMgoUrls()
	mgoApi, err := NewMgoApi(mgoUrls, config.DbName,
		compositePropertyConverter{[]propertyConverter{UUIDToBson, DateToBson}}.convert,
		compositePropertyConverter{[]propertyConverter{UUIDFromBson, DateFromBson, MongoIdRemover}}.convert,
	)
	return mgoApi, err
}

var config, configErr = readConfig()
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
	router.HandleFunc("/{collection}/{resource}", mgoApi.readContent).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", mgoApi.writeContent).Methods("PUT")
	router.HandleFunc("/__health", fthealth.Handler("Dependent services healthceck",
	  "Checking connectivity and usability of dependent services: mongoDB and native-ingester.", mgoHealth))

	http.ListenAndServe(":" + config.Server.Port, nil)
}
