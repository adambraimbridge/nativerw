package resources

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/logging"
	"github.com/gorilla/mux"
)

// Hash hashes the given payload in SHA224 + Hex
func Hash(payload string) string {
	hash := sha256.New224()
	_, err := hash.Write([]byte(payload))
	if err != nil {
		logging.Warn("Failed to write hash!")
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// CheckNativeHash will check for the X-Native-Hash header and compare it to the current saved copy of the same resource
func (f *Filters) CheckNativeHash(mongo db.DB) *Filters {
	next := f.next

	f.next = func(w http.ResponseWriter, r *http.Request) {
		connection, err := mongo.Open()
		if err != nil {
			writeMessage(w, "Failed to connect to the database!", http.StatusServiceUnavailable)
			return
		}

		nativeHash := r.Header.Get("X-Native-Hash")

		if strings.TrimSpace(nativeHash) != "" {
			log := logging.NewTransactionLogger(r.Header.Get(txHeaderKey))
			vars := mux.Vars(r)
			matches, err := checkNativeHash(connection, nativeHash, vars["collection"], vars["resource"])

			if err != nil {
				msg := fmt.Sprintf("Unexpected error occurred while checking the native hash! Message: %v", err.Error())
				log.Error(msg)
				http.Error(w, msg, http.StatusServiceUnavailable)
				return
			}

			if !matches {
				log.Info("The native hash provided with this request does not match the native content in the store, or the original has been removed!")
				http.Error(w, "The native hash provided with this request does not match the native content in the store.", http.StatusConflict)
				return
			}

			writeMessage(w, "Skipped carousel publish.", http.StatusOK)
			return
		}

		next(w, r)
	}

	return f
}

func checkNativeHash(mongo db.Connection, hash string, collection string, id string) (bool, error) {
	resource, found, err := mongo.Read(collection, id)
	if err != nil {
		return false, err
	}

	if !found {
		logging.Warn("Received a carousel publish but the original native content does not exist in the native store! collection=" + collection + ", uuid=" + id)
		return false, nil // no native document for this id, so save it
	}

	data, err := json.Marshal(resource.Content)
	if err != nil {
		return false, err
	}

	existingHash := Hash(string(data))
	return existingHash == hash, nil
}
