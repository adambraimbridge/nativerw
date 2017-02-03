package main

import (
	"os"
	"strings"
	"testing"
)

func startMongo(t *testing.T) *mgoAPI {
	if testing.Short() {
		t.Skip("Mongo integration for long tests only.")
	}

	mongoURL := os.Getenv("MONGO_TEST_URL")
	if strings.TrimSpace(mongoURL) == "" {
		t.Fatal("Please set the environment variable MONGO_TEST_URL to run mongo integration tests (e.g. MONGO_TEST_URL=localhost:27017). Alternatively, run `go test -short` to skip them.")
	}

	conf := configuration{
		Mongos:      mongoURL,
		DbName:      "native-store",
		Collections: []string{"methode"},
	}

	mgo, err := newMgoAPI(&conf)
	if err != nil {
		t.Fatal("Failed to connect to mongo! Please ensure your testing instance is up and running.")
	}

	return mgo
}
