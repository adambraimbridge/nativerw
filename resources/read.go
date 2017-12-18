package resources

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/nativerw/config"
	"github.com/Financial-Times/nativerw/db"
	"github.com/Financial-Times/nativerw/mapper"
	uuidutils "github.com/Financial-Times/uuid-utils-go"
	"github.com/gorilla/mux"
)

// ReadContent reads the native data for the given id and collection
func ReadContent(mongo db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := mongo.Open()
		if err != nil {
			writeMessage(w, "Failed to connect to the database!", http.StatusServiceUnavailable)
			return
		}

		tid := obtainTxID(r)
		vars := mux.Vars(r)
		resourceID := vars["resource"]
		collection := vars["collection"]

		resource, found, err := connection.Read(collection, resourceID)
		if err != nil {
			msg := "Reading from mongoDB failed."
			logger.NewEntry(tid).WithUUID(resourceID).WithError(err).Error(msg)
			http.Error(w, fmt.Sprintf(msg+": %v", err.Error()), http.StatusInternalServerError)
			return
		}

		if !found {
			msg := fmt.Sprintf("Resource not found. collection: %v, id: %v", collection, resourceID)
			logger.NewEntry(tid).WithUUID(resourceID).Info(msg)

			w.Header().Add("Content-Type", "application/json")
			respBody, _ := json.Marshal(map[string]string{"message": msg})
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, string(respBody))
			return
		}

		w.Header().Add("Content-Type", resource.ContentType)

		om, found := mapper.OutMappers[resource.ContentType]
		if !found {
			msg := fmt.Sprintf("Unable to handle resource of type %T", resource)
			logger.NewEntry(tid).WithUUID(resourceID).Warn(msg)
			http.Error(w, msg, http.StatusNotImplemented)
			return
		}

		err = om(w, resource)
		if err != nil {
			msg := fmt.Sprintf("Unable to extract native content from resource with id %v. %v", resourceID, err.Error())
			logger.NewEntry(tid).WithUUID(resourceID).WithError(err).Errorf(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		} else {
			logger.NewEntry(tid).WithUUID(resourceID).Info("Read native content successfully")
		}
	}
}

func ReadIDs(mongo db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := mongo.Open()
		if err != nil {
			writeMessage(w, "Failed to connect to the database!", http.StatusServiceUnavailable)
			return
		}

		vars := mux.Vars(r)
		coll := vars["collection"]
		tid := obtainTxID(r)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		ids, err := connection.ReadIDs(ctx, coll)
		if err != nil {
			msg := fmt.Sprintf(`Failed to read IDs from mongo for %v! "%v"`, coll, err.Error())
			logger.NewEntry(tid).WithError(err).Error(msg)
			http.Error(w, msg, http.StatusServiceUnavailable)
			return
		}

		id := struct {
			ID string `json:"id"`
		}{}

		bw := bufio.NewWriter(w)
		for {
			docID, ok := <-ids
			if !ok {
				break
			}

			id.ID = docID
			jd, _ := json.Marshal(id)

			bw.WriteString(string(jd) + "\n")

			bw.Flush()
			w.(http.Flusher).Flush()
		}
	}
}

func WildcardReadContent(mongo db.DB, conf *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	contentCollections := []string{}
	// loop through each collection that is eligible for wildcard read (i.e. exclude metadata collections)
	for _, c := range conf.Collections {
		if !strings.HasSuffix(c, "metadata") {
			contentCollections = append(contentCollections, c)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		connection, err := mongo.Open()
		if err != nil {
			writeMessage(w, "Failed to connect to the database!", http.StatusServiceUnavailable)
			return
		}

		tid := obtainTxID(r)
		vars := mux.Vars(r)
		resourceID := vars["resource"]
		redirect := isRedirectRequest(r)

		found, collection, res, err := findInCollections(tid, connection, contentCollections, resourceID)
			if found {
				// return 302 with location - or if ?redirect=false, return inline
				location := fmt.Sprintf("../%s/%s", collection, resourceID)
				if redirect {
					w.Header().Add("Location", location)
					writeMessage(w, fmt.Sprintf("See other: %s", location), http.StatusFound)
					//http.Redirect(w, r, location, http.StatusFound)
				} else {
					w.Header().Add("Content-Location", location)
					writeContent(tid, w, res)
				}
				return
			}

		// if still not found - do UUID xor magic
		requestedUUID, _ := uuidutils.NewUUIDFromString(resourceID)
		derivedUUID, _ := uuidutils.NewUUIDDeriverWith(uuidutils.IMAGE_SET).From(requestedUUID)
		derivedID := derivedUUID.String()

		found, collection, res, err = findInCollections(tid, connection, contentCollections, derivedID)
		if found {
			// return 303 with location - or if ?redirect=false, return inline
			location := fmt.Sprintf("../%s/%s", collection, derivedID)
			if redirect {
				w.Header().Add("Location", location)
				writeMessage(w, fmt.Sprintf("See other: %s", location), http.StatusSeeOther)
				//http.Redirect(w, r, location, http.StatusSeeOther)
			} else {
				w.Header().Add("Content-Location", location)
				writeContent(tid, w, res)
			}
			return
		}

		writeMessage(w, fmt.Sprintf("Not found in any collection: %s", resourceID), http.StatusNotFound)
	}
}

func findInCollections(tid string, connection db.Connection, collections []string, uuid string) (bool,string,*mapper.Resource,error) {
	var lastErr error
	for _, c := range collections {
		resource, found, err := connection.Read(c, uuid)
		if err != nil {
			msg := "Reading from mongoDB failed."
			logger.NewEntry(tid).WithUUID(uuid).WithField("collection", c).WithError(err).Error(msg)
			lastErr = err
			continue
		}

		if found {
			logger.NewEntry(tid).WithField("uuid", uuid).WithField("collection", c).Info("Found resource in collection")
			return true, c, &resource, nil
		}

		msg := fmt.Sprintf("Resource not found. collection: %v, id: %v", c, uuid)
		logger.NewEntry(tid).WithUUID(uuid).Info(msg)
	}

	return false, "", nil, lastErr
}

func isRedirectRequest(r *http.Request) bool {
	redirect := r.URL.Query().Get("redirect")
	b, err := strconv.ParseBool(redirect)
	if err != nil {
		b = true
	}

	return b
}

func writeContent(tid string, w http.ResponseWriter, resource *mapper.Resource) {
	w.Header().Add("Content-Type", resource.ContentType)

	om, found := mapper.OutMappers[resource.ContentType]
	if !found {
		msg := fmt.Sprintf("Unable to handle resource of type %T", resource)
		logger.NewEntry(tid).WithUUID(resource.UUID).Warn(msg)
		http.Error(w, msg, http.StatusNotImplemented)
		return
	}

	err := om(w, *resource)
	if err != nil {
		msg := fmt.Sprintf("Unable to extract native content from resource with id %v. %v", resource.UUID, err.Error())
		logger.NewEntry(tid).WithUUID(resource.UUID).WithError(err).Errorf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
	} else {
		logger.NewEntry(tid).WithUUID(resource.UUID).Info("Read native content successfully")
	}
}
