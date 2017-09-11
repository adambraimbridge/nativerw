package resources

import (
	"fmt"
	"github.com/Financial-Times/go-logger"
	"net/http"

	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/mapper"
	"github.com/gorilla/mux"
)

// WriteContent writes a new native record
func WriteContent(mongo db.DB) func(w http.ResponseWriter, r *http.Request) {
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

		contentType := r.Header.Get("Content-Type")
		inMapper := mapper.InMappers[contentType]
		if inMapper == nil {
			msg := fmt.Sprintf("Content-Type header missing. Default value ('application/octet-stream') is used.")
			logger.NewEntry(tid).WithUUID(resourceID).Info(msg)
			contentType = "application/octet-stream"
			inMapper = mapper.InMappers[contentType]
		}

		content, err := inMapper(r.Body)
		if err != nil {
			msg := "Extracting content from HTTP body failed"
			logger.NewMonitoringEntry("NativeSave", tid, "").WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusBadRequest)
			return
		}

		wrappedContent := mapper.Wrap(content, resourceID, contentType)

		if err := connection.Write(collectionID, wrappedContent); err != nil {
			msg := "Writing to mongoDB failed"
			logger.NewMonitoringEntry("NativeSave", tid, "").WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusInternalServerError)
			return
		}

		logger.NewMonitoringEntry("NativeSave", tid, "").WithUUID(resourceID).Info("Successfully saved")
	}
}
