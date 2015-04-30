package main

import (
    "log"
    "os"
)

type CombinedLogger struct {
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
}

func (l CombinedLogger) info(s string) {
    l.Info.Println(s);
}

func (l CombinedLogger) warn(s string) {
    l.Warning.Println(s);
}

func (l CombinedLogger) error(s string) {
    l.Error.Println(s);
}

var logger CombinedLogger

func initLoggers() {
    logger = CombinedLogger {
        log.New(os.Stdout, "INFO    - ", log.Ldate|log.Ltime|log.Lshortfile),
        log.New(os.Stdout, "WARNING - ", log.Ldate|log.Ltime|log.Lshortfile),
        log.New(os.Stderr, "ERROR   - ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}
