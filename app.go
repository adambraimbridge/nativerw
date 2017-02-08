package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Financial-Times/nativerw/config"
	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/logging"
	"github.com/Financial-Times/nativerw/resources"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/kr/pretty"
)

func main() {
	cliApp := cli.App("nativerw", "Writes any raw content/data from native CMS in mongoDB without transformation.")
	mongos := cliApp.String(cli.StringOpt{
		Name:   "mongos",
		Value:  "",
		Desc:   "Mongo addresses to connect to in format: host1[:port1][,host2[:port2],...]",
		EnvVar: "MONGOS",
	})

	configFile := cliApp.String(cli.StringOpt{
		Name:   "config",
		Value:  "config.json",
		Desc:   "Config file (e.g. config.json)",
		EnvVar: "CONFIG",
	})

	cliApp.Action = func() {
		logging.Info("Starting nativerw app.")

		conf, err := config.ReadConfig(*configFile)
		if err != nil {
			logging.Error(fmt.Sprintf("Error reading the configuration: %+v\n", err.Error()))
			os.Exit(1)
		}

		if len(*mongos) != 0 {
			conf.Mongos = *mongos
		}

		logging.Info(fmt.Sprintf("Using configuration %# v \n", pretty.Formatter(conf)))

		mongo, err := db.NewDatabase(conf)
		for err != nil {
			logging.Error(fmt.Sprintf("Couldn't establish connection to mongoDB: %+v", err.Error()))

			time.Sleep(5 * time.Second)

			mongo, err = db.NewDatabase(conf)
		}

		logging.Info("Established connection to mongoDB.")
		mongo.EnsureIndex()

		router(mongo)
		err = http.ListenAndServe(":"+conf.Server.Port, nil)
		if err != nil {
			logging.Error(fmt.Sprintf("Couldn't set up HTTP listener: %+v\n", err))
		}
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		println(err)
	}
}

func router(mongo db.DB) *mux.Router {
	router := mux.NewRouter()
	http.HandleFunc("/", resources.AccessLogging(router))
	router.HandleFunc("/{collection}/__ids", resources.Filter(resources.GetIDs(mongo)).ValidateAccessForCollection(mongo).Build()).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", resources.Filter(resources.ReadContent(mongo)).ValidateAccess(mongo).Build()).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", resources.Filter(resources.WriteContent(mongo)).ValidateAccess(mongo).CheckNativeHash(mongo).Build()).Methods("PUT")
	router.HandleFunc("/{collection}/{resource}", resources.DeleteContent(mongo)).Methods("DELETE")
	router.HandleFunc("/__health", resources.Healthchecks(mongo))
	router.HandleFunc("/__gtg", resources.GoodToGo(mongo))
	return router
}