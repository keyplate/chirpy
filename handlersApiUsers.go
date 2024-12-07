package main

import (
    "encoding/json"
    "io"
    "net/http"
    "time"

    "github.com/keyplate/chirpy/internal/database"
    "github.com/google/uuid"
)

type createUserBody struct {
    Email string `json:"email"`
}

type usrResponse struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    dat, err := io.ReadAll(req.Body)
    var usrEmail createUserBody

    if err != nil {
        respondWithError(w, 400, "Can not read body")
        return
    }

    err = json.Unmarshal(dat, &usrEmail)
    if err != nil {    
        respondWithError(w, 400, "Can not read body")
        return
    }
    
    usrParams := database.CreateUserParams{ ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Email: usrEmail.Email }
    usr, err := cfg.db.CreateUser(req.Context(), usrParams)
    if err != nil {
        respondWithError(w, 500, err.Error())
        return
    }

    usrResp := usrResponse{ ID: usr.ID, CreatedAt: usr.CreatedAt, UpdatedAt: usr.UpdatedAt, Email: usr.Email }
    respondWithJSON(w, 201, usrResp)
}
