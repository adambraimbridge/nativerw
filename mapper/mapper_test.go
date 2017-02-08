package mapper

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	var tests = []struct {
		resource         map[string]interface{}
		uuid             string
		contentType      string
		publishReference string
		wantResource     Resource
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
		result := Wrap(test.resource, test.uuid, test.contentType)
		if !reflect.DeepEqual(result, test.wantResource) {
			t.Errorf("Resource: %v\n, Expected: %v\n, Actual: %v", test.resource, test.wantResource, result)
		}
	}
}

func TestJsonMappers(t *testing.T) {
	testResource := Resource{
		UUID: "9694733e-163a-4393-801f-000ab7de5041",
		Content: map[string]interface{}{
			"title":  "Title",
			"body":   "This is a body.",
			"brands": []interface{}{"Lex", "Markets"},
		},
		ContentType: "application/json",
	}

	var writer = bytes.NewBuffer([]byte{})

	outMapper := OutMappers["application/json"]
	err := outMapper(writer, testResource)

	assert.NoError(t, err, "Shouldn't error")
	assert.Equal(t, `{"body":"This is a body.","brands":["Lex","Markets"],"title":"Title"}`, strings.TrimSpace(writer.String()), "Json should match")

	inMapper := InMappers["application/json"]
	actual, err := inMapper(strings.NewReader(`{"body":"This is a body.","brands":["Lex","Markets"],"title":"Title"}`))

	assert.NoError(t, err)
	for k, v := range actual.(map[string]interface{}) {
		assert.EqualValues(t, testResource.Content.(map[string]interface{})[k], v)
	}
}

func TestBinaryMappers(t *testing.T) {
	testResource := Resource{
		UUID:        "9694733e-163a-4393-801f-000ab7de5041",
		Content:     []byte("hi"),
		ContentType: "application/json",
	}

	var writer = bytes.NewBuffer([]byte{})

	outMapper := OutMappers["application/octet-stream"]
	err := outMapper(writer, testResource)

	assert.NoError(t, err, "Shouldn't error")
	assert.Equal(t, `hi`, strings.TrimSpace(writer.String()))

	inMapper := InMappers["application/octet-stream"]
	actual, err := inMapper(strings.NewReader(`hi`))

	assert.NoError(t, err)
	assert.Equal(t, testResource.Content, actual)
}
