package main

import (
	"reflect"
	"testing"
	"errors"
)

func TestValidateAccess(t *testing.T) {

	config := Configuration{
		Collections:[]string{"methode", "wordpress"},
		Mongos: "localhost:27017",
		DbName: "dbname",
		Server: Server{Port:"port"},
	}
	api, _ := NewMgoApi(&config)

	var tests = []struct {
		collectionId  string
		resourceId    string
		expectedError error
	}{
		{
			"methode",
			"9694733e-163a-4393-801f-000ab7de5041",
			nil,
		},
		{
			"wordpress",
			"9694733e-163a-4393-801f-000ab7de5041",
			nil,
		},
		{
			"other",
			"9694733e-163a-4393-801f-000ab7de5041",
			errors.New("Collection not supported or resourceId not a valid uuid."),
		},
	}

	for _, test := range tests {
		actual_err := api.validateAccess(test.collectionId, test.resourceId)
		if (actual_err != test.expectedError && actual_err.Error() != test.expectedError.Error()) {
			t.Errorf("Expected: %v\n, Actual: %v", test.expectedError, actual_err)
		}
	}
}

func TestWrap(t *testing.T) {

	var tests = []struct {
		resource     map[string]interface{}
		uuid         string
		contentType  string
		wantResource Resource
	}{
		{
			map[string]interface{}{
				"title":  "Title",
				"body":   "This is a body.",
				"brands": []string{"Lex", "Markets"},
			},
			"9694733e-163a-4393-801f-000ab7de5041",
			"application/json",
			Resource{
				UUID: "9694733e-163a-4393-801f-000ab7de5041",
				Content: map[string]interface{}{
					"title":  "Title",
					"body":   "This is a body.",
					"brands": []string{"Lex", "Markets"},
				},
				ContentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		result := wrap(test.resource, test.uuid, test.contentType)
		if !reflect.DeepEqual(result, test.wantResource) {
			t.Errorf("Resource: %v\n, Expected: %v\n, Actual: %v", test.resource, test.wantResource, result)
		}
	}
}
