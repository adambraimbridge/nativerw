package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (s3a *s3api) getIds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coll := vars["collection"]
	enc := json.NewEncoder(w)
	stop := make(chan struct{})
	defer close(stop)
	all, err := s3a.Ids(coll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := struct {
		ID string `json:"id"`
	}{}
	for _, docID := range all {
		id.ID = docID
		enc.Encode(id)
	}
}

func (s3a *s3api) readContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	uuid := vars["uuid"]

	content, err := s3a.Read(collection, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", content.contentType)
	_, err = io.Copy(w, bytes.NewReader(content.content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s3a *s3api) writeContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	uuid := vars["uuid"]
	body, _ := ioutil.ReadAll(r.Body)
	contentType := r.Header.Get("Content-Type")

	resource := &s3resource{
		uuid:        uuid,
		content:     body,
		contentType: contentType,
	}

	err := s3a.Write(collection, *resource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s3a *s3api) deleteContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	collection := vars["collection"]
	uuid := vars["uuid"]

	err := s3a.Delete(collection, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
