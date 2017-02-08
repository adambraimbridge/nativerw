package resources

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetIDs(t *testing.T) {
	mongo := new(MockDB)
	channel := make(chan string)
	mongo.On("Ids", "methode", mock.AnythingOfType("chan struct {}"), mock.AnythingOfType("chan error")).Return(channel)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/__ids", GetIDs(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/__ids", nil)

	go func() {
		time.Sleep(250 * time.Millisecond)
		channel <- "i am a real id"
		time.Sleep(10 * time.Millisecond)
		close(channel)
	}()

	router.ServeHTTP(w, req)

	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody := strings.Split(w.Body.String(), "\n")
	assert.Equal(t, `{"id":"i am a real id"}`, strings.TrimSpace(responseBody[0]))
}

func TestGetIDsFail(t *testing.T) {
	mongo := new(MockDB)
	channel := make(chan string)
	mongo.On("Ids", "methode", mock.AnythingOfType("chan struct {}"), mock.AnythingOfType("chan error")).Return(channel)

	router := mux.NewRouter()
	router.HandleFunc("/{collection}/__ids", GetIDs(mongo)).Methods("GET")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/methode/__ids", nil)

	go func() {
		time.Sleep(250 * time.Millisecond)

		channel <- "i am a real id"

		stop := mongo.CallArgs[2].(chan error)
		stop <- errors.New("stop now pls")
	}()

	router.ServeHTTP(w, req)

	mongo.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, w.Code)

	responseBody := strings.Split(w.Body.String(), "\n")

	assert.Equal(t, `{"id":"i am a real id"}`, strings.TrimSpace(responseBody[0]))
	assert.Equal(t, `stop now pls`, strings.TrimSpace(responseBody[1]))
}
