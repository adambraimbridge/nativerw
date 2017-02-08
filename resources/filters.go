package resources

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/logging"
	"github.com/gorilla/mux"
)

var uuidRegexp = regexp.MustCompile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")

func validateAccess(mongo db.DB, collectionID, resourceID string) error {
	if mongo.GetSupportedCollections()[collectionID] && uuidRegexp.MatchString(resourceID) {
		return nil
	}
	return errors.New("Collection not supported or resourceId not a valid uuid.")
}

func validateAccessForCollection(mongo db.DB, collectionID string) error {
	if mongo.GetSupportedCollections()[collectionID] {
		return nil
	}
	return errors.New("Collection not supported.")
}

// ValidateAccess validates whether the collection exists and the resource ID is in uuid format.
func (f *Filters) ValidateAccess(mongo db.DB) *Filters {
	next := f.next
	f.next = func(w http.ResponseWriter, r *http.Request) {
		collectionID := mux.Vars(r)["collection"]
		resourceID := mux.Vars(r)["resource"]

		if err := validateAccess(mongo, collectionID, resourceID); err != nil {
			ctxlogger := logging.NewTransactionLogger(obtainTxID(r))
			msg := fmt.Sprintf("Invalid collectionId (%v) or resourceId (%v).\n%v", collectionID, resourceID, err)
			ctxlogger.Info(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		next(w, r)
	}
	return f
}

// ValidateAccessForCollection validates whether the collection exists
func (f *Filters) ValidateAccessForCollection(mongo db.DB) *Filters {
	next := f.next
	f.next = func(w http.ResponseWriter, r *http.Request) {
		collection := mux.Vars(r)["collection"]

		if err := validateAccessForCollection(mongo, collection); err != nil {
			ctxLogger := logging.NewTransactionLogger(obtainTxID(r))
			msg := fmt.Sprintf("Invalid collectionId (%v).\n%v", collection, err)
			ctxLogger.Info(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		next(w, r)
	}
	return f
}

// Filters wraps the next http handler
type Filters struct {
	next func(w http.ResponseWriter, r *http.Request)
}

// Filter creates a new composable filter.
func Filter(next func(w http.ResponseWriter, r *http.Request)) *Filters {
	return &Filters{next}
}

// Build returns the final chained handler
func (f *Filters) Build() func(w http.ResponseWriter, r *http.Request) {
	return f.next
}
