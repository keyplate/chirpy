package main

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
)

var censoredWords = []string{"kerfuffle", "sharbert", "fornax"}

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
    
    if !validateChirpLength(chirp.Body) {
        respondWithError(w, 400, "Chirp is too long")
        return 
    }
    
    cleanedChirp := validateAndReplaceProfane(chirp.Body)
    respondWithJSON(w, 200, map[string]string{"cleaned_body": cleanedChirp})
}

func validateChirpLength(chirp string) bool {
    if len(chirp) > 140 {
        return false
    }
    return true
}

func validateAndReplaceProfane(chirp string) string {
    words := strings.Split(chirp, " ")
    for i, word := range(words) {
        if slices.Contains(censoredWords, strings.ToLower(word)) {
            words[i] = "****"
        }
    }
    return strings.Join(words, " ")
}
