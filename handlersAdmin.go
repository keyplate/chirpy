package main

import (
    "fmt"
    "net/http"
    "io"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    hits := cfg.reqCount.Load()
    io.WriteString(w, fmt.Sprintf(
    `<html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
      </body>
    </html>`,
    hits))
    w.WriteHeader(http.StatusOK)
}

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
