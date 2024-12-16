package main

import (
	"net/http"
	"time"

	"github.com/keyplate/chirpy/internal/auth"
)

type tokenResponse struct {
    Token string `json:"token"`
}

func (cfg *apiConfig)handlerRefersh(w http.ResponseWriter, req *http.Request) {
    refreshTokenStr, err := auth.GetBearerToken(req.Header)
    if err != nil {
        respondWithError(w, 401, "Unauthorized") 
        return
    }

    refreshToken, err := cfg.db.GetRefreshToken(req.Context(), refreshTokenStr)
    if err != nil || refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
        respondWithError(w, 401, "Unauthorized") 
        return
    }

    usr, err := cfg.db.GetUserByToken(req.Context(), refreshToken.Token)
    if err != nil {
        respondWithError(w, 401, err.Error())
        return
    }

    token, err := auth.MakeJWT(usr.ID, cfg.secret, initExpirationDur(3600))
    if err != nil {
        respondWithError(w, 401, "Unauthorized")
    }

    respondWithJSON(w, 200, tokenResponse{ Token: token })
}
