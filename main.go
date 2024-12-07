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
}

func main() {
    godotenv.Load()

    connStr := os.Getenv("DB_URL") 
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    dbQueries := database.New(db)

    serveMux := http.NewServeMux()
    cfg := apiConfig{ reqCount: atomic.Int32{}, db: dbQueries }

    appHandler :=  http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
    serveMux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))
    serveMux.Handle("GET /api/healthz", http.HandlerFunc(handlerHeatlthz))
    serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
    serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)
    serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

    server := http.Server{ Handler: serveMux, Addr: ":8080" }
    err = server.ListenAndServe()
    if err != nil {
        fmt.Printf("Error: %v", err) 
    }
}
