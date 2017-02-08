package resources

import (
	"fmt"
	"net/http"

	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/logging"
	"github.com/gorilla/mux"
)

// DeleteContent deletes the given resource from the given collection
func DeleteContent(mongo db.DB) func(writer http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		collectionID := mux.Vars(r)["collection"]
		resourceID := mux.Vars(r)["resource"]
		ctxlogger := logging.NewTransactionLogger(obtainTxID(r))

		if err := mongo.Delete(collectionID, resourceID); err != nil {
			msg := fmt.Sprintf("Deleting from mongoDB failed:\n%v\n", err)
			ctxlogger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		ctxlogger.Info(fmt.Sprintf("Delete native content successful. resource_id: %+v", resourceID))
	}
}
