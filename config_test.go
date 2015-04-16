package main

import (
	"testing"
	"strings"
)

func TestReadConfigFromReader(t *testing.T) {

	var tests = [] struct {
		json, wantDbName, wantPort string
	} {

		{"{ \"dbName\": \"ft\", \"server\": { \"port\" : \"1234\" } }", "ft", "1234"},

	}

	for _, c := range tests {
		got, _ := readConfigFromReader(strings.NewReader(c.json))
		var resDbName = got.DbName
		var resServerPort = got.Server.Port
		if resDbName != c.wantDbName ||
			 resServerPort != c.wantPort {
			t.Errorf("Input Json: %q, resultDbName: %q and resultPort: %q", c.json, resDbName, resServerPort)
		}
	}

}

func TestPrepareMgoUrls(t *testing.T) {
	var tests = [] struct {
		config Configuration
		wantUrl string
	} {
		{
			Configuration {
				[]Mongo {
					Mongo {
						"localhost",
						"1000",
					},
					Mongo {
						"localhost",
						"1001",
					},
				},
				"testdb",
				Server {
					"9999",
				},
			},
			"localhost:1000,localhost:1001",
		},
	}

	for _, test := range tests {
		result := test.config.prepareMgoUrls()
		if result != test.wantUrl {
			t.Errorf("\nMongos: %v\nExpected: %v\nActual: %v", test.config, test.wantUrl, result)
		}
	}
}