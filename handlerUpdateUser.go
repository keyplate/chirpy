package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/keyplate/chirpy/internal/auth"
	"github.com/keyplate/chirpy/internal/database"
)

type updateUsrRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type updateUsrResponse struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
    IsChirpyRed bool `json:"is_chirpy_red"`
} 

func (cfg *apiConfig)handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    token, err := auth.GetBearerToken(req.Header) 
    if err != nil {
        respondWithError(w, 401, "Unauthorized")
        return
    }

    usrID, err := auth.ValidateJWT(token, cfg.secret)
    if err != nil {
        respondWithError(w, 401, "Unauthorized")
        return
    }

    var updateUsrReq updateUsrRequest
    dat, err := io.ReadAll(req.Body)
    if err != nil {
        respondWithError(w, 500, "Something went worng")
        return
    }

    err = json.Unmarshal(dat, &updateUsrReq)
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    }

    hashedPassword, err := auth.HashPassword(updateUsrReq.Password)
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    }

    updatedUsr, err := cfg.db.UpdateUserEmailPassword(req.Context(), database.UpdateUserEmailPasswordParams{
        ID: usrID,  
        Email: updateUsrReq.Email,
        HashedPassword: hashedPassword,
    })
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
    }

    respondWithJSON(w, 200, updateUsrResponse{
        ID: updatedUsr.ID,
        CreatedAt: updatedUsr.CreatedAt,
        UpdatedAt: updatedUsr.UpdatedAt,
        Email: updatedUsr.Email,
        IsChirpyRed: updatedUsr.IsChirpyRed,
    })
}
