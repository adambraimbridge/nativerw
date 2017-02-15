package resources

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Financial-Times/nativerw/logging"
)

const txHeaderKey = "X-Request-Id"
const txHeaderLength = 20

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// AccessLogging intercepts traffic and logs the request and response
func AccessLogging(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		ctxlogger := logging.NewTransactionLogger(obtainTxID(r))
		var loggingWriter = &logging.ResponseLogger{Writer: w}

		next.ServeHTTP(loggingWriter, r)

		ctxlogger.Access(fmt.Sprintf("status=%v ; method=%v ; url=%v ; response_time=%.4f ; response_size=%v",
			loggingWriter.Status(),
			r.Method,
			r.URL.String(),
			time.Now().Sub(t1).Seconds(),
			loggingWriter.Size()))
	}
}

func writeMessage(w http.ResponseWriter, msg string, status int) {
	data, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{msg})

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(status)
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
