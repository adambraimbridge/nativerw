package main

import (
	"net/http"
	"os"
	"strconv"

	logger "github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/nativerw/config"
	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/resources"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"github.com/kr/pretty"
)

const appName = "nativerw"

func init() {
	logger.InitLogger(appName, "info")
}

func main() {
	cliApp := cli.App("nativerw", "Writes any raw content/data from native CMS in mongoDB without transformation.")
	mongos := cliApp.String(cli.StringOpt{
		Name:   "mongos",
		Value:  "",
		Desc:   "Mongo addresses to connect to in format: host1:port1[,host2:port2,...]",
		EnvVar: "MONGOS",
	})

	mongoNodeCount := cliApp.Int(cli.IntOpt{
		Name:   "mongo_node_count",
		Value:  3,
		Desc:   "Number of mongoDB instances",
		EnvVar: "MONGO_NODE_COUNT",
	})

	configFile := cliApp.String(cli.StringOpt{
		Name:   "config",
		Value:  "config.json",
		Desc:   "Config file (e.g. config.json)",
		EnvVar: "CONFIG",
	})

	cliApp.Action = func() {
		conf, err := config.ReadConfig(*configFile)
		if err != nil {
			logger.Fatalf(nil, err, "Error reading the configuration")
		}

		if err := db.CheckMongoUrls(*mongos, *mongoNodeCount); err != nil {
			logger.Fatalf(nil, err, "Provided mongoDB urls %s are invalid", *mongos)
		}

		conf.Mongos = *mongos
		logger.Infof(nil, "Using configuration %# v", pretty.Formatter(conf))

		logger.ServiceStartedEvent(conf.Server.Port)
		mongo := db.NewDBConnection(conf)
		router(mongo)

		go func() {
			connection, err := mongo.Open()
			if err != nil {
				logger.Errorf(nil, err, "Mongo connection not yet established, awaiting stable connection")
				connection, err = mongo.Await()
				if err != nil {
					logger.Fatalf(nil, err, "Unrecoverable error connecting to mongo")
				}
			}

			logger.Infof(map[string]interface{}{}, "Established connection to mongoDB.")
			connection.EnsureIndex()
		}()

		err = http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), nil)
		if err != nil {
			logger.Fatalf(nil, err, "Couldn't set up HTTP listener")
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

	router.HandleFunc("/{collection}/__ids", resources.Filter(resources.ReadIDs(mongo)).ValidateAccessForCollection(mongo).Build()).Methods("GET")

	router.HandleFunc("/{collection}/{resource}", resources.Filter(resources.ReadContent(mongo)).ValidateAccess(mongo).Build()).Methods("GET")
	router.HandleFunc("/{collection}/{resource}", resources.Filter(resources.WriteContent(mongo)).ValidateAccess(mongo).CheckNativeHash(mongo).Build()).Methods("PUT")
	router.HandleFunc("/{collection}/{resource}", resources.Filter(resources.DeleteContent(mongo)).ValidateAccess(mongo).Build()).Methods("DELETE")

	router.HandleFunc("/__health", resources.Healthchecks(mongo))
	router.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(resources.GoodToGo(mongo)))

	router.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler).Methods("GET")
	router.HandleFunc(status.PingPath, status.PingHandler).Methods("GET")

	return router
}
