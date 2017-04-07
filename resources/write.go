package resources

import (
	"fmt"
	"net/http"

	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/logging"
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
		ctxlogger := logging.NewTransactionLogger(obtainTxID(r))

		contentType := r.Header.Get("Content-Type")
		inMapper := mapper.InMappers[contentType]
		if inMapper == nil {
			msg := fmt.Sprintf("Content-Type header missing. Default value ('application/octet-stream') is used.")
			ctxlogger.Info(msg)
			contentType = "application/octet-stream"
			inMapper = mapper.InMappers[contentType]
		}

		content, err := inMapper(r.Body)
		if err != nil {
			msg := fmt.Sprintf("Extracting content from HTTP body failed:\n%v\n", err)
			ctxlogger.Warn(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		wrappedContent := mapper.Wrap(content, resourceID, contentType)

		if err := connection.Write(collectionID, wrappedContent); err != nil {
			msg := fmt.Sprintf("Writing to mongoDB failed:\n%v\n", err)
			ctxlogger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		ctxlogger.Info(fmt.Sprintf("Written native content. resource_id: %+v", resourceID))
	}
}
