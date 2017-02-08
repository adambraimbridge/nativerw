package resources

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/logging"
	"github.com/Financial-Times/nativerw/mapper"
	"github.com/gorilla/mux"
)

// ReadContent reads the native data for the given id and collection
func ReadContent(mongo db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		resourceID := vars["resource"]
		collection := vars["collection"]
		ctxlogger := logging.NewTransactionLogger(obtainTxID(r))

		found, resource, err := mongo.Read(collection, resourceID)
		if err != nil {
			msg := fmt.Sprintf("Reading from mongoDB failed.\n%v\n", err.Error())
			ctxlogger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if !found {
			msg := fmt.Sprintf("Resource not found. collection: %v, id: %v", collection, resourceID)
			ctxlogger.Warn(msg)

			w.Header().Add("Content-Type", "application/json")
			respBody, _ := json.Marshal(map[string]string{"message": msg})
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(respBody))
			return
		}

		w.Header().Add("Content-Type", resource.ContentType)

		om, found := mapper.OutMappers[resource.ContentType]
		if !found {
			msg := fmt.Sprintf("Unable to handle resource of type %T. resourceId: %v, resource: %v", resource, resourceID, resource)
			ctxlogger.Warn(msg)
			http.Error(w, msg, http.StatusNotImplemented)
			return
		}

		err = om(w, resource)
		if err != nil {
			msg := fmt.Sprintf("Unable to extract native content from resource with id %v. %v", resourceID, err.Error())
			ctxlogger.Warn(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		} else {
			ctxlogger.Info(fmt.Sprintf("Read native content. resource_id: %+v", resourceID))
		}
	}
}
