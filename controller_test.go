package main

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		resource         map[string]interface{}
		uuid             string
		contentType      string
		publishReference string
		wantResource     resource
	}{
		{
			map[string]interface{}{
				"title":  "Title",
				"body":   "This is a body.",
				"brands": []string{"Lex", "Markets"},
			},
			"9694733e-163a-4393-801f-000ab7de5041",
			"application/json",
			"tid_blahblahblah",
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
		result := wrap(test.resource, test.uuid, test.contentType, test.publishReference)
		if !reflect.DeepEqual(result, test.wantResource) {
			t.Errorf("Resource: %v\n, Expected: %v\n, Actual: %v", test.resource, test.wantResource, result)
		}
	}
}

func TestJsonOutMapper(t *testing.T) {
	testResource := resource{
		UUID: "9694733e-163a-4393-801f-000ab7de5041",
		Content: map[string]interface{}{
			"title":  "Title",
			"body":   "This is a body.",
			"brands": []string{"Lex", "Markets"},
		},
		ContentType: "application/json",
	}

	var writer = bytes.NewBuffer([]byte{})

	jsonMapper := outMappers["application/json"]
	err := jsonMapper(writer, testResource)

	assert.NoError(t, err, "Shouldn't error")
	assert.Equal(t, `{"body":"This is a body.","brands":["Lex","Markets"],"title":"Title"}`, strings.TrimSpace(writer.String()), "Json should match")
}

func TestWrite(t *testing.T) {
	initLoggers()
	mongo := startMongo(t)
	defer mongo.session.Close()

	router := router(mongo)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/methode/163ccf41-0134-4abc-95cc-7d419591edd6", strings.NewReader(`{"uuid":"163ccf41-0134-4abc-95cc-7d419591edd6", "title": "Donald: In His Own Words"}`))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Request-Id", "tid_my-fake-tid")

	router.ServeHTTP(resp, req)

	found, res, err := mongo.Read("methode", "163ccf41-0134-4abc-95cc-7d419591edd6")

	assert.True(t, found, "Should be found")
	assert.NoError(t, err, "Should not error")
	assert.Equal(t, "application/json", res.ContentType, "Should match")
	assert.Equal(t, "163ccf41-0134-4abc-95cc-7d419591edd6", res.UUID, "Should match")
}
