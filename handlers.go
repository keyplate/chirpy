package main

import (
    "fmt"
    "io"
	"net/http"
)

func handlerHeatlthz(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte("OK"))
    w.WriteHeader(http.StatusOK)
}

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
    cfg.reqCount.Store(0)
    w.WriteHeader(http.StatusOK)
}
