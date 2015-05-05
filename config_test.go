package main

import (
	"testing"
)

func TestPrepareMgoUrls(t *testing.T) {
	var tests = []struct {
		config  Configuration
		wantUrl string
	}{
		{
			Configuration{
				[]string{
					"localhost:1000",
					"localhost:1001",
				},
				"testdb",
				Server{
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
