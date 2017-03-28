package resources

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestObtainTxID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/doesnt/matter", nil)
	req.Header.Add("X-Request-Id", "tid_blahblah")
	txid := obtainTxID(req)
	assert.Equal(t, "tid_blahblah", txid)
}

func TestObtainTxIDGeneratesANewOneIfNoneAvailable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/doesnt/matter", nil)
	txid := obtainTxID(req)
	assert.Contains(t, txid, "tid_")
}

func TestAccessLogWillForwardRequest(t *testing.T) {
	passed := false
	next := func(w http.ResponseWriter, r *http.Request) {
		passed = true
	}

	router := mux.NewRouter()
	router.HandleFunc("/fake", next).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/fake", nil)

	AccessLogging(router)(w, req)
	assert.True(t, passed)
}
