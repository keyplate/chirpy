package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/keyplate/chirpy/internal/auth"
	"github.com/keyplate/chirpy/internal/database"
)

func (cfg *apiConfig)handlerRevoke(w http.ResponseWriter, req *http.Request) {
    refreshTokenStr, err := auth.GetBearerToken(req.Header)
    if err != nil {
        respondWithError(w, 500, "Something went wrong")
        return
    } 

    _, err = cfg.db.MarkRevoked(
        req.Context(),
        database.MarkRevokedParams{
            UpdatedAt: time.Now(),
            RevokedAt: sql.NullTime{
                Time: time.Now(),
                Valid: true,
            },
            Token: refreshTokenStr,
        },
    )
    if err != nil {
        respondWithError(w, 500, "Something went wrong")  
        return
    }

    respondWithJSON(w, 204, "")
}
