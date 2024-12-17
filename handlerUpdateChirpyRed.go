package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/keyplate/chirpy/internal/auth"
)

type userUpgradeRequest struct {
    Event string `json:"event"`
    Data struct { 
        UserID string `json:"user_id"`
    } `json:"data"`
}

var eventType string = "user.upgraded"

func (cfg *apiConfig)handlerUpdateChirpyRed(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()
    
    token, err := auth.GetAPIKey(req.Header)
    if err != nil || token != cfg.polkaAPIKey {
        respondWithError(w, 401, "Unauthorized")
        return
    }

    dat, err := io.ReadAll(req.Body)
    if err != nil {
        respondWithError(w, 400, "Couldn't read body")
        return
    }

    var upgradeRequest userUpgradeRequest
    err = json.Unmarshal(dat, &upgradeRequest)
    if err != nil {
       respondWithError(w, 400, "Couldn't process body")
       return
    }

    if upgradeRequest.Event != eventType {
        respondWithJSON(w, 204, "")
        return
    }

    usrID, err := uuid.Parse(upgradeRequest.Data.UserID)
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    }
    
    _, err = cfg.db.UpdateIsChirpyRedUserTrue(req.Context(), usrID)
    if err != nil {
        respondWithError(w, 404, "Not Found")
        return
    }

    respondWithJSON(w, 204, "")
}
