package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/joho/godotenv"
	"github.com/keyplate/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
    db *database.Queries
    reqCount atomic.Int32
    platform string
    secret string
    polkaAPIKey string
}

func main() {
    godotenv.Load()

    connStr := os.Getenv("DB_URL")
    platfrom := os.Getenv("PLATFORM")
    secret := os.Getenv("JWT_SECRET")
    polkaAPIKey := os.Getenv("POLKA_KEY")
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    dbQueries := database.New(db)

    serveMux := http.NewServeMux()
    cfg := apiConfig{ 
        reqCount: atomic.Int32{}, 
        db: dbQueries, 
        platform: platfrom, 
        secret: secret, 
        polkaAPIKey: polkaAPIKey,
}

    appHandler :=  http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
    serveMux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))
    serveMux.Handle("GET /api/healthz", http.HandlerFunc(handlerHeatlthz))
    
    serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
    serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
    
    serveMux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
    serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
    serveMux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
    serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)

    serveMux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
    serveMux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
    serveMux.HandleFunc("POST /api/login", cfg.handlerLogin)
    
    serveMux.HandleFunc("POST /api/refresh", cfg.handlerRefersh)
    serveMux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

    serveMux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpdateChirpyRed)

    server := http.Server{ Handler: serveMux, Addr: ":8080" }
    err = server.ListenAndServe()
    if err != nil {
        fmt.Printf("Error: %v", err) 
    }
}
