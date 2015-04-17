package main

import (
	"fmt"
	"git.svc.ft.com/scm/gl/fthealth.git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func createMgoApi(config *Configuration) (*MgoApi, error) {
	mgoUrls := config.prepareMgoUrls()
	mgoApi, err := NewMgoApi(mgoUrls, config.DbName,
		compositePropertyConverter{[]propertyConverter{UUIDToBson, DateToBson}}.convert,
		compositePropertyConverter{[]propertyConverter{UUIDFromBson, DateFromBson, MongoIdRemover}}.convert,
	)
	return mgoApi, err
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Missing parameter. Usage: <pathToExecutable>/nativerw <confFilePath>\n")
		return
	}

	config, configErr := readConfig(os.Args[1])
	if configErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", configErr.Error())
		return
	}

	mgoApi, mgoApiCreationErr := createMgoApi(config)
	if mgoApiCreationErr != nil {
		fmt.Fprintf(os.Stderr, "%v\n", mgoApiCreationErr.Error())
		return
	}

	router := mux.NewRouter()
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, router))
	router.HandleFunc("/{collection}/{resource}", mgoApi.readContent).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", mgoApi.writeContent).Methods("PUT")
	router.HandleFunc("/__health", fthealth.Handler("Dependent services healthceck",
		"Checking connectivity and usability of dependent services: mongoDB and native-ingester.", mgoApi.buildHealthCheck()))

	http.ListenAndServe(":"+config.Server.Port, nil)
}
