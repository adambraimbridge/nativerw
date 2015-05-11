package main

import (
	"fmt"
	"git.svc.ft.com/scm/gl/fthealth.git"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	initLoggers()
	logger.info("Starting nativerw app.")

	if len(os.Args) < 2 {
		logger.error("Missing parameter. Usage: <pathToExecutable>/nativerw <pathToConfigurationFile>\n")
		return
	}

	config, configErr := readConfig(os.Args[1])
	if configErr != nil {
		logger.error(fmt.Sprintf("Error reading the configuration: %+v\n", configErr.Error()))
		return
	}

	mgoApi, mgoApiCreationErr := NewMgoApi(config)
	if mgoApiCreationErr != nil {
		logger.error(fmt.Sprintf("Couldn't establish connection to mongoDB: %+v\n", mgoApiCreationErr.Error()))
		return
	}
	mgoApi.EnsureIndex()

	router := mux.NewRouter()
	http.Handle("/", handlers.CombinedLoggingHandler(logger, router))
	router.HandleFunc("/{collection}/{resource}", mgoApi.readContent).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", mgoApi.writeContent).Methods("PUT")
	router.HandleFunc("/__health", fthealth.Handler("Dependent services healthcheck",
		"Checking connectivity and usability of dependent services: mongoDB.",
		mgoApi.writeHealthCheck(), mgoApi.readHealthCheck()))
	router.HandleFunc("/__gtg", mgoApi.goodToGo)
	router.HandleFunc("/__test/log/info", logDummyInfo).Methods("POST")
	router.HandleFunc("/__test/log/warning", logDummyWarn).Methods("POST")
	router.HandleFunc("/__test/log/error", logDummyError).Methods("POST")
	err := http.ListenAndServe(":"+config.Server.Port, nil)
	if err != nil {
		logger.error(fmt.Sprintf("Couldn't set up HTTP listener: %+v\n", err))
	}
}
