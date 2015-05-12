package main

import (
	"net/http"
	"time"
	"fmt"
)

type accessLoggingHandler struct {
	underlyingHandler http.Handler
}

func (h accessLoggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t1 := time.Now()
	ctxlogger := TxCombinedLogger{logger, obtainTxId(req)}
	var loggingWriter = &responseLogger{w: w}

	h.underlyingHandler.ServeHTTP(loggingWriter, req)

	ctxlogger.access(fmt.Sprintf("status=%v ; method=%v ; url=%v ; response_time=%.4f ; response_size=%v",
		loggingWriter.Status(),
		req.Method,
		req.URL.String(),
		time.Now().Sub(t1).Seconds(),
		loggingWriter.Size()))
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	if l.status == 0 {
		return http.StatusOK
	}
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}
