package main


import (
	"net/http"

	"github.com/google/uuid"
	"github.com/keyplate/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
    token, err := auth.GetBearerToken(req.Header)
    if err != nil {
        respondWithError(w, 401, "Unauthenticated")
        return 
    }

    usrID, err := auth.ValidateJWT(token, cfg.secret)
    if err != nil {
        respondWithError(w, 401, "Unauthenticated")
        return
    }

    chirpID, err := uuid.Parse(req.PathValue("chirpID"))
    if err != nil {
        respondWithError(w, 400, "Couldn't parse chirp uuid")
        return
    }


    chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
    if err != nil {
        respondWithError(w, 404, "Not Found")
        return
    }

    if chirp.UserID != usrID {
        respondWithError(w, 403, "Forbidden")
        return
    }

    err = cfg.db.DeleteChirpByID(req.Context(), chirp.ID)
    if err != nil {
        respondWithError(w, 500, "Something went worng")
        return
    }

    respondWithJSON(w, 204, "No Content")
}
