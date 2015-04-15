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
			t.Errorf("for Json: %q, resultDbName is %q and resultPort is %q", c.json, resDbName, resServerPort)
		}
	}

}