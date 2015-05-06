package main

import (
	"fmt"
	"log"
	"os"
)

type CombinedLogger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Access  *log.Logger
}

func (l CombinedLogger) info(txId, msg string) {
	l.Info.Println(fmt.Sprintf("transaction_id: %+v - %+v", txId, msg))
}

func (l CombinedLogger) warn(txId, msg string) {
	l.Warning.Println(fmt.Sprintf("transaction_id: %+v - %+v", txId, msg))
}

func (l CombinedLogger) error(txId, msg string) {
	l.Error.Println(fmt.Sprintf("transaction_id: %+v - %+v", txId, msg))
}

func (l CombinedLogger) access(msg string) {
	l.Access.Println(msg)
}

var logger CombinedLogger

func initLoggers() {
	logger = CombinedLogger{
		log.New(os.Stdout, "INFO    - ", log.Ldate|log.Ltime|log.Lshortfile),
		log.New(os.Stdout, "WARNING - ", log.Ldate|log.Ltime|log.Lshortfile),
		log.New(os.Stderr, "ERROR   - ", log.Ldate|log.Ltime|log.Lshortfile),
		log.New(os.Stdout, "ACCESS  - ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

type AccessWriter struct {
	logger *log.Logger
}

func (w AccessWriter) Write(p []byte) (int, error) {
	msg := string(p[:])
	w.logger.Print(msg)
	return len(msg), nil
}
