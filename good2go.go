package main

import "net/http"


func (ma *MgoApi) goodToGo(writer http.ResponseWriter, req *http.Request) {
    healthChecks := [] func() error { ma.checkReadable, ma.checkWritable };

    for _, hCheck := range healthChecks{
        if err := hCheck(); err != nil {
            writer.WriteHeader(http.StatusServiceUnavailable)
            return
        }
    }
}