package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

// func main() {
// 	cliApp := cli.App("nativerw", "Writes any raw content/data from native CMS in mongoDB without transformation.")
// 	mongos := cliApp.String(cli.StringOpt{
// 		Name:   "mongos",
// 		Value:  "",
// 		Desc:   "Mongo addresses to connect to in format: host1[:port1][,host2[:port2],...]",
// 		EnvVar: "MONGOS",
// 	})
// 	config := cliApp.String(cli.StringOpt{
// 		Name:   "config",
// 		Value:  "config.json",
// 		Desc:   "Config file (e.g. config.json)",
// 		EnvVar: "CONFIG",
// 	})
// 	cliApp.Action = func() {
// 		initLoggers()
// 		logger.info("Starting nativerw app.")
// 		config, configErr := readConfig(*config)
// 		if configErr != nil {
// 			logger.error(fmt.Sprintf("Error reading the configuration: %+v\n", configErr.Error()))
// 			os.Exit(1)
// 		}
// 		if len(*mongos) != 0 {
// 			config.Mongos = *mongos
// 		}
//
// 		logger.info(fmt.Sprintf("Using configuration %# v \n", pretty.Formatter(config)))
// 		mgoAPI, mgoAPICreationErr := newMgoAPI(config)
// 		for mgoAPICreationErr != nil {
// 			logger.error(fmt.Sprintf("Couldn't establish connection to mongoDB: %+v", mgoAPICreationErr.Error()))
// 			time.Sleep(5 * time.Second)
// 			mgoAPI, mgoAPICreationErr = newMgoAPI(config)
// 		}
// 		logger.info("Established connection to mongoDB.")
// 		mgoAPI.EnsureIndex()
//
// 		router := mux.NewRouter()
// 		http.Handle("/", accessLoggingHandler{router})
// 		router.HandleFunc("/{collection}/__ids", mgoAPI.getIds).Methods("GET")
// 		router.HandleFunc("/{collection}/{resource}", mgoAPI.readContent).Methods("GET")
// 		router.HandleFunc("/{collection}/{resource}", mgoAPI.writeContent).Methods("PUT")
// 		router.HandleFunc("/__health", fthealth.Handler("Dependent services healthcheck",
// 			"Checking connectivity and usability of dependent services: mongoDB.",
// 			mgoAPI.writeHealthCheck(), mgoAPI.readHealthCheck()))
// 		router.HandleFunc("/__gtg", mgoAPI.goodToGo)
// 		err := http.ListenAndServe(":"+config.Server.Port, nil)
// 		if err != nil {
// 			logger.error(fmt.Sprintf("Couldn't set up HTTP listener: %+v\n", err))
// 		}
// 	}
// 	err := cliApp.Run(os.Args)
// 	if err != nil {
// 		println(err)
// 	}
// }

func main() {
	cliApp := cli.App("nativerw", "Writes any raw content/data from native CMS into S3 without transformation.")
	s3Bucket := cliApp.String(cli.StringOpt{
		Name:   "s3bucket",
		Value:  "com.ft.coco-native-store.semantic",
		Desc:   "S3 bucket of the native store content",
		EnvVar: "S3_BUCKET",
	})

	cliApp.Action = func() {
		initLoggers()
		logger.info("Starting nativerw app")

		s3api, err := newS3API(*s3Bucket)
		for err != nil {
			logger.error("Error in connecting to S3")
			return
		}

		router := mux.NewRouter()
		http.Handle("/", accessLoggingHandler{router})
		router.HandleFunc("/{collection}/__ids", s3api.getIds).Methods("GET")
		router.HandleFunc("/{collection}/{uuid}", s3api.readContent).Methods("GET")
		router.HandleFunc("/{collection}/{uuid}", s3api.writeContent).Methods("PUT")
		router.HandleFunc("/{collection}/{uuid}", s3api.deleteContent).Methods("DELETE")
		err = http.ListenAndServe(":8082", nil)
		if err != nil {
			logger.error(fmt.Sprintf("Couldn't set up HTTP listener: %+v\n", err))
		}
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		println(err)
	}
}
