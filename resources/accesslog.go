package resources

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

const txHeaderKey = "X-Request-Id"
const txHeaderLength = 20

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// AccessLogging intercepts traffic and logs the request and response
func AccessLogging(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}

func writeMessage(w http.ResponseWriter, msg string, status int) {
	data, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{msg})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

func obtainTxID(req *http.Request) string {
	txID := req.Header.Get(txHeaderKey)
	if txID == "" {
		return "tid_" + randSeq(txHeaderLength)
	}
	return txID
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
