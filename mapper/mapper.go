package mapper

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Resource is the representation of a native resource
type Resource struct {
	UUID        string
	Content     interface{}
	ContentType string
}

// Wrap creates a new resource
func Wrap(content interface{}, resourceID, contentType string) Resource {
	return Resource{
		UUID:        resourceID,
		Content:     content,
		ContentType: contentType,
	}
}

// OutMapper writes a resource in the required content format
type OutMapper func(io.Writer, Resource) error

// OutMappers contains all the supported mappers
var OutMappers = map[string]OutMapper{
	"application/json": func(w io.Writer, resource Resource) error {
		encoder := json.NewEncoder(w)
		return encoder.Encode(resource.Content)
	},
	"application/octet-stream": func(w io.Writer, resource Resource) error {
		data := resource.Content.([]byte)
		_, err := io.Copy(w, bytes.NewReader(data))
		return err
	},
}

// InMapper marshals the transport format into a resource
type InMapper func(io.Reader) (interface{}, error)

// InMappers contains all the supported mappers
var InMappers = map[string]InMapper{
	"application/json": func(r io.Reader) (interface{}, error) {
		var c map[string]interface{}
		err := json.NewDecoder(r).Decode(&c)
		return c, err
	},
	"application/octet-stream": func(r io.Reader) (interface{}, error) {
		return ioutil.ReadAll(r)
	},
}
