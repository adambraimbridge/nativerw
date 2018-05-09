package resources

import (
	"fmt"
	"net/http"
	"strings"

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

		contentTypeHeader := r.Header.Get("Content-Type")
		contentType := contentTypeHeader

		// in case, the content-type header comes with additional parameters
		if strings.Contains(contentType, ";") {
			contentType = strings.Split(contentType, ";")[0]
		}

		inMapper := mapper.InMappers[contentType]
		if inMapper == nil {
			msg := fmt.Sprintf("Content-Type header missing. Default value ('application/octet-stream') is used.")
			logger.WithTransactionID(tid).WithUUID(resourceID).Warn(msg)
			contentTypeHeader = "application/octet-stream"
			inMapper = mapper.InMappers[contentTypeHeader]
		}

		content, err := inMapper(r.Body)
		if err != nil {
			msg := "Extracting content from HTTP body failed"
			logger.WithMonitoringEvent("NativeSave", tid, "").WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusBadRequest)
			return
		}

		wrappedContent := mapper.Wrap(content, resourceID, contentTypeHeader)

		if err := connection.Write(collectionID, wrappedContent); err != nil {
			msg := "Writing to mongoDB failed"
			logger.WithMonitoringEvent("NativeSave", tid, "").WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+":\n%v\n", err), http.StatusInternalServerError)
			return
		}

		logger.WithMonitoringEvent("NativeSave", tid, "").WithUUID(resourceID).Info("Successfully saved")
	}
}
