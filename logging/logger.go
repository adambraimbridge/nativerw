package logging

import (
	"fmt"
	"log"
	"os"
)

// DefaultLogger provides access to the default logger
func DefaultLogger() CombinedLogger {
	return logger
}

// Info convenience method for writing to the default logger info
func Info(msg string) {
	logger.Info(msg)
}

// Warn convenience method for writing to the default logger warn
func Warn(msg string) {
	logger.Warn(msg)
}

// Error convenience method for writing to the default logger error
func Error(msg string) {
	logger.Error(msg)
}

// Access convenience method for writing to the default logger access
func Access(msg string) {
	logger.Access(msg)
}

// CombinedLogger the general logger interface
type CombinedLogger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Access(msg string)
	Write(p []byte) (int, error)
}

// SimpleCombinedLogger is the basic underlying implementation of the logger.
type SimpleCombinedLogger struct {
	info    *log.Logger
	warning *log.Logger
	err     *log.Logger
	access  *log.Logger
}

// Info implementation of the info func
func (l SimpleCombinedLogger) Info(msg string) {
	l.info.Println(msg)
}

// Warn implementation of the warn func
func (l SimpleCombinedLogger) Warn(msg string) {
	l.warning.Println(msg)
}

// Error implementation of the error func
func (l SimpleCombinedLogger) Error(msg string) {
	l.err.Println(msg)
}

// Access implementation of the access func
func (l SimpleCombinedLogger) Access(msg string) {
	l.access.Println(msg)
}

func (l SimpleCombinedLogger) Write(p []byte) (int, error) {
	msg := string(p)
	l.access.Print(msg)
	return len(msg), nil
}

var logger CombinedLogger

func init() {
	logger = SimpleCombinedLogger{
		info:    log.New(os.Stdout, "INFO    - ", log.Ldate|log.Ltime),
		warning: log.New(os.Stdout, "WARNING - ", log.Ldate|log.Ltime),
		err:     log.New(os.Stderr, "ERROR   - ", log.Ldate|log.Ltime),
		access:  log.New(os.Stdout, "ACCESS  - ", log.Ldate|log.Ltime),
	}
}

// TxCombinedLogger wraps the underlying logger, and appends transaction_id to the logline.
type TxCombinedLogger struct {
	Wrapped       CombinedLogger
	TransactionID string
}

// NewTransactionLogger returns a new TxID aware logger.
func NewTransactionLogger(txid string) CombinedLogger {
	return TxCombinedLogger{logger, txid}
}

// Info implementation of the info func
func (l TxCombinedLogger) Info(msg string) {
	l.Wrapped.Info(fmt.Sprintf("transaction_id=%+v - %+v", l.TransactionID, msg))
}

// Warn implementation of the warn func
func (l TxCombinedLogger) Warn(msg string) {
	l.Wrapped.Warn(fmt.Sprintf("transaction_id=%+v - %+v", l.TransactionID, msg))
}

// Error implementation of the error func
func (l TxCombinedLogger) Error(msg string) {
	l.Wrapped.Error(fmt.Sprintf("transaction_id=%+v - %+v", l.TransactionID, msg))
}

// Access implementation of the access func
func (l TxCombinedLogger) Access(msg string) {
	l.Wrapped.Access(fmt.Sprintf("transaction_id=%+v - %+v", l.TransactionID, msg))
}

func (l TxCombinedLogger) Write(p []byte) (int, error) {
	msg := string(p)
	l.Access(msg)
	return len(msg), nil
}
