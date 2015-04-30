package main

import (
    "log"
    "os"
    "io/ioutil"
    "io"
)

type CombinedLogger struct {
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger

    Access  io.Writer
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
        ioutil.Discard,
    }
}

func initAccessLog(config *Configuration) error {
    accessFile, err := os.OpenFile(config.Server.AccessLogs, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        return err
    }
//    defer accessFile.Close()
    logger.Access = accessFile
    return nil
}
