package resources

import (
	"encoding/json"
	"net/http"

	"github.com/Financial-Times/nativerw/db"
	"github.com/gorilla/mux"
)

// GetIDs writes all ids from the native collection
func GetIDs(mongo db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		coll := vars["collection"]

		enc := json.NewEncoder(w)
		stop := make(chan struct{})
		errChan := make(chan error)

		defer close(stop)
		defer close(errChan)

		all := mongo.Ids(coll, stop, errChan)
		id := struct {
			ID string `json:"id"`
		}{}

		for {
			select {
			case err := <-errChan:
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			case docID, ok := <-all:
				if !ok {
					return
				}
				id.ID = docID
				enc.Encode(id)
			}
		}
	}
}
