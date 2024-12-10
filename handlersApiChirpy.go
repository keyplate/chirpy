package main

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/keyplate/chirpy/internal/database"
)

var censoredWords = []string{"kerfuffle", "sharbert", "fornax"}

type chirpRequest struct {
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body string `json:"body"`
    UserID uuid.UUID `json:"user_id"`
}

func handlerHeatlthz(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte("OK"))
    w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig)handlerGetChirps(w http.ResponseWriter, req *http.Request) {
    chirps, err := cfg.db.GetAllChirps(req.Context())
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    }

    chirpResponseList := []chirpResponse{}
    for _, chirp := range(chirps) {
        chirpResponseList = append(chirpResponseList, toChirpResponse(chirp))
    }
    respondWithJSON(w, 200, chirpResponseList)
}

func (cfg *apiConfig)handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {
    chirpID, err := uuid.Parse(req.PathValue("chirpID"))
    if err != nil {
        respondWithError(w, 400, "Chirp ID not valid")
        return
    }
    
    chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
    if err != nil {
        respondWithError(w, 404, "Chirp not found")
        return
    }
    
    respondWithJSON(w, 200, toChirpResponse(chirp))
}

func (cfg *apiConfig)handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    var chirpRequest chirpRequest
    data, err := io.ReadAll(req.Body)
    if err != nil {
        respondWithError(w, 400, "Can not read body")
        return
    }
    
    err = json.Unmarshal(data, &chirpRequest)
    if err != nil {
        respondWithError(w, 400, "Can not decode body")
        return
    }
    
    if !validateChirpLength(chirpRequest.Body) {
        respondWithError(w, 400, "Chirp is too long")
        return 
    }
    
    cleanedChirp := validateAndReplaceProfane(chirpRequest.Body)

    createChirpParams := database.CreateChirpParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Body: cleanedChirp,
        UserID: chirpRequest.UserID,
    }
    
    chirp, err := cfg.db.CreateChirp(req.Context(), createChirpParams)
    if err != nil {
        respondWithError(w, 400, "Can not create chirp")
    }

    respondWithJSON(w, 201, toChirpResponse(chirp))
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

func toChirpResponse(chirp database.Chirp) chirpResponse {
    return chirpResponse{
        ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        Body: chirp.Body,
        UserID: chirp.UserID,
    }
}
