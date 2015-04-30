package main

import (
	"fmt"
	"git.svc.ft.com/scm/cp/fthealth-go.git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var logger CombinedLogger

func createMgoApi(config *Configuration) (*MgoApi, error) {
	mgoUrls := config.prepareMgoUrls()
	mgoApi, err := NewMgoApi(mgoUrls, config.DbName)
	return mgoApi, err
}

func main() {
    initLoggers()
    logger.info("Starting nativerw app.")

	if len(os.Args) < 2 {
        logger.error("Missing parameter. Usage: <pathToExecutable>/nativerw <pathToConfigurationFile>\n")
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
	router.HandleFunc("/healthcheck", fthealth.Handler("Dependent services healthceck",
		"Checking connectivity and usability of dependent services: mongoDB.",
		mgoApi.writeHealthCheck(), mgoApi.readHealthCheck()))
	err := http.ListenAndServe(":"+config.Server.Port, nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
