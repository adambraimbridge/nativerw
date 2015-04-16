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
