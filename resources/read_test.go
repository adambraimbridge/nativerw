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

func TestReadContent(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Read", "methode", "a-real-uuid").Return(true, mapper.Resource{ContentType: "application/json", Content: map[string]interface{}{"uuid": "fake-data"}}, nil)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", ReadContent(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/a-real-uuid", nil)

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"uuid":"fake-data"}`, strings.TrimSpace(w.Body.String()))
}

func TestReadFailed(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Read", "methode", "a-real-uuid").Return(false, mapper.Resource{}, errors.New("i failed"))

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", ReadContent(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/a-real-uuid", nil)

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestIDNotFound(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Read", "methode", "a-real-uuid").Return(false, mapper.Resource{}, nil)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", ReadContent(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/a-real-uuid", nil)

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNoMapperImplemented(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Read", "methode", "a-real-uuid").Return(true, mapper.Resource{ContentType: "application/vnd.fake-mime-type"}, nil)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", ReadContent(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/a-real-uuid", nil)

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusNotImplemented, w.Code)
}

func TestUnableToMap(t *testing.T) {
	mongo := new(MockDB)
	mongo.On("Read", "methode", "a-real-uuid").Return(true, mapper.Resource{ContentType: "application/json", Content: func() {}}, nil)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/{resource}", ReadContent(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/a-real-uuid", nil)

	router.ServeHTTP(w, req)
	mongo.AssertExpectations(t)
	t.Log(w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}