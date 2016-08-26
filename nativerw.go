package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth"
	"github.com/gorilla/mux"
)

func main() {
	initLoggers()
	logger.info("Starting nativerw app.")

	mongos := flag.String("mongos", "", "Mongo addresses to connect to in format: host1[:port1][,host2[:port2],...]")
	flag.Parse()

	if len(flag.Args()) == 0 {
		logger.error("Missing parameter. Usage: <pathToExecutable>/nativerw <pathToConfigurationFile>\n")
		os.Exit(1)
	}

	config, configErr := readConfig(flag.Arg(0))
	if configErr != nil {
		logger.error(fmt.Sprintf("Error reading the configuration: %+v\n", configErr.Error()))
		os.Exit(1)
	}
	if len(*mongos) != 0 {
		config.Mongos = *mongos
	}

	mgoAPI, mgoAPICreationErr := newMgoAPI(config)
	for {
		if mgoAPICreationErr != nil {
			logger.error(fmt.Sprintf("Couldn't establish connection to mongoDB: %+v", mgoAPICreationErr.Error()))
			time.Sleep(5 * time.Second)
			mgoAPI, mgoAPICreationErr = newMgoAPI(config)
		} else {
			logger.info("Established connection to mongoDB.")
			break
		}
	}

	mgoAPI.EnsureIndex()

	router := mux.NewRouter()
	http.Handle("/", accessLoggingHandler{router})
	router.HandleFunc("/{collection}/__ids", mgoAPI.getIds).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", mgoAPI.readContent).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", mgoAPI.writeContent).Methods("PUT")
	router.HandleFunc("/__health", fthealth.Handler("Dependent services healthcheck",
		"Checking connectivity and usability of dependent services: mongoDB.",
		mgoAPI.writeHealthCheck(), mgoAPI.readHealthCheck()))
	router.HandleFunc("/__gtg", mgoAPI.goodToGo)
	router.HandleFunc("/__test/log/info", logDummyInfo).Methods("POST")
	router.HandleFunc("/__test/log/warning", logDummyWarn).Methods("POST")
	router.HandleFunc("/__test/log/error", logDummyError).Methods("POST")
	router.HandleFunc("/__test/slow", testSlowRequest).Methods("POST")
	err := http.ListenAndServe(":"+config.Server.Port, nil)
	if err != nil {
		logger.error(fmt.Sprintf("Couldn't set up HTTP listener: %+v\n", err))
	}
}
