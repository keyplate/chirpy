package main

import (
    "encoding/json"
    "io"
    "net/http"
    "time"

    "github.com/keyplate/chirpy/internal/auth"
    "github.com/keyplate/chirpy/internal/database"
    "github.com/google/uuid"
)

type createUserRequest struct {
    Email string `json:"email"`
    Password string `json:"password"` 
}

type loginUserRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type usrResponse struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    var usrReq createUserRequest
    dat, err := io.ReadAll(req.Body)
    if err != nil {
        respondWithError(w, 400, "Can not read body")
        return
    }

    err = json.Unmarshal(dat, &usrReq)
    if err != nil {    
        respondWithError(w, 400, "Can not read body")
        return
    }
    
    hashedPass, err := auth.HashPassword(usrReq.Password)
    if err != nil {
        respondWithError(w, 400, "Something went wrong")
    }
    usrParams := database.CreateUserParams{ 
        ID: uuid.New(), 
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Email: usrReq.Email,
        HashedPassword: hashedPass,
    }
    usr, err := cfg.db.CreateUser(req.Context(), usrParams)
    if err != nil {
        respondWithError(w, 500, err.Error())
        return
    }

    usrResp := usrResponse{ ID: usr.ID, CreatedAt: usr.CreatedAt, UpdatedAt: usr.UpdatedAt, Email: usr.Email }
    respondWithJSON(w, 201, usrResp)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()


    var loginReq loginUserRequest
    dat, err := io.ReadAll(req.Body)
    if err != nil {
        respondWithError(w, 400, "Can not read body")
        return
    }

    err = json.Unmarshal(dat, &loginReq)
    if err != nil {
        respondWithError(w, 400, "Can not read body")
        return
    }

    usr, err := cfg.db.GetUserByEmail(req.Context(), loginReq.Email)
    if err != nil {
        respondWithError(w, 401, "Unauthorized")
        return
    }

    err = auth.CheckPasswordHash(loginReq.Password, usr.HashedPassword)
    if err != nil {
        respondWithError(w, 401, "Unauthorized")
        return
    }

    respondWithJSON(w, 200, usrResponse{
        ID: usr.ID,
        CreatedAt: usr.CreatedAt,
        UpdatedAt: usr.UpdatedAt,
        Email: usr.Email,
    })
}
