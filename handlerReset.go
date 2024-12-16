package main

import (
    "fmt"
    "net/http"
    "io"
)



func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
    if cfg.platform != "dev" {
        respondWithError(w, 403, "Forbidden")
        return
    }

    err := cfg.db.DeleteAllUsers(req.Context())
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    }
    cfg.reqCount.Store(0)
    w.WriteHeader(http.StatusOK)
}
