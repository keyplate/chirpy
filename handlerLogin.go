package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/keyplate/chirpy/internal/auth"
	"github.com/keyplate/chirpy/internal/database"
)

type loginUsrRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type loginUsrResponse struct {
    ID uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email string `json:"email"`
    IsChirpyRed bool `json:"is_chirpy_red"`
    Token string `json:"token"`
    RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
    defer req.Body.Close()

    var loginReq loginUsrRequest
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

    expiresIn := initExpirationDur(3600)
    token, err := auth.MakeJWT(usr.ID, cfg.secret, expiresIn)
    if err != nil {
        respondWithError(w, 401, "Unauthorized")
        return
    }

    refreshToken, err := cfg.issueRefreshToken(usr.ID, req.Context())
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    }

    respondWithJSON(w, 200, loginUsrResponse{
        ID: usr.ID,
        CreatedAt: usr.CreatedAt,
        UpdatedAt: usr.UpdatedAt,
        Email: usr.Email,
        Token: token,
        RefreshToken: refreshToken.Token,
        IsChirpyRed: usr.IsChirpyRed,
    })
}

func initExpirationDur(seconds int) time.Duration {
    hourInSec := 3600
    if seconds <= 0 || seconds > hourInSec {
        return time.Duration(float64(hourInSec) * float64(time.Second))
    }
    return time.Duration(float64(seconds) * float64(time.Second))
}

func (cfg *apiConfig)issueRefreshToken(usrID uuid.UUID, ctx context.Context) (database.RefreshToken, error) {
    revokeDuration, err := time.ParseDuration("1440h") //60 days
    if err != nil {
        return database.RefreshToken{}, err
    }
    revokeAt := time.Now().Add(revokeDuration)

    refrethTokenStr, err := auth.MakeRefreshToken()
    if err != nil {
        return database.RefreshToken{}, err 
    }
    
    createRefreshTokenParams := database.CreateRefreshTokenParams{
        Token: refrethTokenStr,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: usrID,
        ExpiresAt: revokeAt,
        RevokedAt: sql.NullTime{},
    }

    refreshToken, err := cfg.db.CreateRefreshToken(ctx, createRefreshTokenParams)
    
    return refreshToken, nil
}
