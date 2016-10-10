package main

import (
	"encoding/json"
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
