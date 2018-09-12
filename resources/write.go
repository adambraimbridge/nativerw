package resources

import (
	"fmt"
	"net/http"

	"github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/mapper"
	"github.com/gorilla/mux"
)

// WriteContent writes a new native record
func WriteContent(mongo db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		connection, err := mongo.Open()
		if err != nil {
			writeMessage(w, "Failed to connect to the database!", http.StatusServiceUnavailable)
			return
		}

		collectionID := mux.Vars(r)["collection"]
		resourceID := mux.Vars(r)["resource"]
		tid := obtainTxID(r)

		contentTypeHeader := extractAttrFromHeader(r, "Content-Type", "application/octet-stream", tid, resourceID)

		inMapper, err := mapper.InMapperForContentType(contentTypeHeader)
		if err != nil {
			msg := "Unsupported content-type"
			logger.WithMonitoringEvent("NativeSave", tid, contentTypeHeader).WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusBadRequest)
			return
		}

		originSystemIDHeader := extractAttrFromHeader(r, "Origin-System-Id", "", tid, resourceID)
		content, err := inMapper(r.Body)
		if err != nil {
			msg := "Extracting content from HTTP body failed"
			logger.WithMonitoringEvent("NativeSave", tid, "").WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusBadRequest)
			return
		}

		wrappedContent := mapper.Wrap(content, resourceID, contentTypeHeader, originSystemIDHeader)

		if err := connection.Write(collectionID, wrappedContent); err != nil {
			msg := "Writing to mongoDB failed"
			logger.WithMonitoringEvent("NativeSave", tid, "").WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusInternalServerError)
			return
		}

		logger.WithMonitoringEvent("NativeSave", tid, "").WithUUID(resourceID).Info("Successfully saved")
	}
}
