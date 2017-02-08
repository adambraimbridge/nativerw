package resources

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Financial-Times/nativerw/mapper"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestWriteContent(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Write", "methode", mapper.Resource{UUID: "a-real-uuid", Content: map[string]interface{}{}, ContentType: "application/json"}).Return(nil)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", WriteContent(mongo)).Methods("PUT")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/methode/a-real-uuid", strings.NewReader(`{}`))

	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWriteFailed(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Write", "methode", mapper.Resource{UUID: "a-real-uuid", Content: map[string]interface{}{}, ContentType: "application/json"}).Return(errors.New("i failed"))

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", WriteContent(mongo)).Methods("PUT")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/methode/a-real-uuid", strings.NewReader(`{}`))

	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDefaultsToBinaryMapping(t *testing.T) {
	mongo := new(MockDB)
	content, _ := mapper.InMappers["application/octet-stream"](strings.NewReader(`{}`))

	mongo.On("Write", "methode", mapper.Resource{UUID: "a-real-uuid", Content: content, ContentType: "application/octet-stream"}).Return(errors.New("i failed"))

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", WriteContent(mongo)).Methods("PUT")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/methode/a-real-uuid", strings.NewReader(`{}`))

	req.Header.Add("Content-Type", "application/a-fake-type")

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestFailedJSON(t *testing.T) {
	mongo := new(MockDB)
	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", WriteContent(mongo)).Methods("PUT")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/methode/a-real-uuid", strings.NewReader(`i am not json`))

	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
