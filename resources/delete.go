package resources

import (
	"fmt"
	"github.com/Financial-Times/go-logger"
	"net/http"

	"github.com/Financial-Times/nativerw/db"
	"github.com/gorilla/mux"
)

// DeleteContent deletes the given resource from the given collection
func DeleteContent(mongo db.DB) func(writer http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := mongo.Open()
		if err != nil {
			writeMessage(w, "Failed to connect to the database!", http.StatusServiceUnavailable)
			return
		}

		defer r.Body.Close()
		collectionID := mux.Vars(r)["collection"]
		resourceID := mux.Vars(r)["resource"]
		tid := obtainTxID(r)

		if err := connection.Delete(collectionID, resourceID); err != nil {
			msg := "Deleting from mongoDB failed"
			logger.ErrorEventWithUUID(tid, resourceID, msg, err)
			http.Error(w, fmt.Sprintf(msg+"\n%v\n", err), http.StatusInternalServerError)
			return
		}

		logger.MonitoringEventWithUUID("NativeDelete", tid, resourceID, "Annotations", "Successfully deleted")
	}
}
