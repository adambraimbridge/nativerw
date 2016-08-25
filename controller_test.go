package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestValidateAccess(t *testing.T) {

	collectionMappings := createMapWithAllowedCollections([]string{"methode", "wordpress"})
	api := &mgoAPI{"", nil, collectionMappings}

	var tests = []struct {
		collectionID  string
		resourceID    string
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
		actualError := api.validateAccess(test.collectionID, test.resourceID)
		if actualError != test.expectedError && actualError.Error() != test.expectedError.Error() {
			t.Errorf("Expected: %v\n, Actual: %v", test.expectedError, actualError)
		}
	}
}

func TestWrap(t *testing.T) {

	var tests = []struct {
		resource     map[string]interface{}
		uuid         string
		contentType  string
		wantResource resource
	}{
		{
			map[string]interface{}{
				"title":  "Title",
				"body":   "This is a body.",
				"brands": []string{"Lex", "Markets"},
			},
			"9694733e-163a-4393-801f-000ab7de5041",
			"application/json",
			resource{
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
