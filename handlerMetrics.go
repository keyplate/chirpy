package main

import (
    "net/http"
    "io"
    "fmt"
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
