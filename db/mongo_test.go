package db

import (
	"os"
	"strings"
	"testing"

	"github.com/Financial-Times/nativerw/config"
)

func startMongo(t *testing.T) DB {
	if testing.Short() {
		t.Skip("Mongo integration for long tests only.")
	}

	mongoURL := os.Getenv("MONGO_TEST_URL")
	if strings.TrimSpace(mongoURL) == "" {
		t.Fatal("Please set the environment variable MONGO_TEST_URL to run mongo integration tests (e.g. export MONGO_TEST_URL=localhost:27017). Alternatively, run `go test -short` to skip them.")
	}

	conf := config.Configuration{
		Mongos:      mongoURL,
		DbName:      "native-store",
		Collections: []string{"methode"},
	}

	mgo, err := NewDBConnection(&conf)
	if err != nil {
		t.Fatal("Failed to connect to mongo! Please ensure your testing instance is up and running.")
	}

	return mgo
}
