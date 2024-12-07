package main

import (
    "encoding/json"
    "fmt"
    "io"
	"net/http"
)

type chirpBody struct {
    Body string `json:"body"`
}

func handlerHeatlthz(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte("OK"))
    w.WriteHeader(http.StatusOK)
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
    body := req.Body
    var chirp chirpBody

    data, err := io.ReadAll(body)
    if err != nil {
        respondWithError(w, 400, "Can not read body")
        return
    }
    
    err = json.Unmarshal(data, &chirp)
    if err != nil {
        respondWithError(w, 400, "Can not decode body")
        return
    }
    
    if !validateChirp(chirp.Body) {
        respondWithError(w, 400, "Chirp is too long")
        return 
    }
    respondWithJSON(w, 200, map[string]bool{"valid": true})
}

func validateChirp(chirp string) bool {
    if len(chirp) > 140 {
        return false
    }
    return true
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
