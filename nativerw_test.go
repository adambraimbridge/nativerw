package main

import (
"testing"
"strings"
)

func TestWrapResource(t *testing.T) {

	var tests = [] struct {
		resource map[string]interface, uuid, contentType string,
		wantResource map[string]interface
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
				"content": {
					"title": "Title", 
					"body": "This is a body.",
					"brands" : []string {"Lex", "Markets"},
				},
				"content-type" : "application/json",
			}
		}
	}

	for _, test := range tests {
		result := wrapResource(test.resource, test.uuid, test.contentType)
		if result != test.wantResource {
			t.Errorf("Resource: %v, Expected: %v, Actual: %v", test.resource, test.wantResource, result)
		}
	}
}
