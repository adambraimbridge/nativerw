package main

import (
	"testing"
	"reflect"
)

func TestWrapResource(t *testing.T) {

	var tests = [] struct {
		resource map[string]interface{}
		uuid string
		contentType string
		wantResource map[string]interface{}
	} {
		{ 			
			map[string]interface{} {
				"title": "Title", 
				"body": "This is a body.",
				"brands" : []string {"Lex", "Markets"},
			},
			"9694733e-163a-4393-801f-000ab7de5041",
			"application/json",
			map[string]interface{} {
				"uuid": "9694733e-163a-4393-801f-000ab7de5041",
				"content": map[string]interface{} {
					"title": "Title", 
					"body": "This is a body.",
					"brands" : []string {"Lex", "Markets"},
				},
				"content-type" : "application/json",
			},
		},
	}

	for _, test := range tests {
		result := wrapMap(test.resource, test.uuid, test.contentType)
		if !reflect.DeepEqual(result, test.wantResource) {
			t.Errorf("Resource: %v\n, Expected: %v\n, Actual: %v", test.resource, test.wantResource, result)
		}
	}
}

func TestUnwrapResource(t *testing.T) {
	var tests = [] struct {
		wrappedResource map[string]interface{}
		wantResource map[string]interface{}
	} {
		{
			map[string]interface{} {
				"uuid": "9694733e-163a-4393-801f-000ab7de5041",
				"content": map[string]interface{} {
					"title": "Title",
					"body": "This is a body.",
					"brands" : []string {"Lex", "Markets"},
				},
				"content-type" : "application/json",
			},
			map[string]interface{} {
				"title": "Title",
				"body": "This is a body.",
				"brands" : []string {"Lex", "Markets"},
			},
		},
	}

	for _, test := range tests {
		result := unwrapResource(test.wrappedResource)
		if !reflect.DeepEqual(result, test.wantResource) {
			t.Errorf("Resource: %v\n, Expected: %v\n, Actual: %v", test.wrappedResource, test.wantResource, result)
		}
	}	
}

func TestPrepareMgoUrls(t *testing.T) {
	var tests = [] struct {
		mongos []Mongo
		wantUrl string
	} {
		{
			[]Mongo{
				{"localhost", "1000"},
				{"localhost", "1001"},
			},
			"localhost:1000,localhost:1001",
		},
	}

	for _, test := range tests {
		result := prepareMgoUrls(test.mongos)
		if result != test.wantUrl {
			t.Errorf("\nMongos: %v\nExpected: %v\nActual: %v", test.mongos, test.wantUrl, result)
		}
	}

}