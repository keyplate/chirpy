package main

import (
    "fmt"
    "net/http"
    "sync/atomic"
)

type apiConfig struct {
    reqCount atomic.Int32 
}

func main() {
    serveMux := http.NewServeMux()
    cfg := apiConfig{ reqCount: atomic.Int32{} }

    appHandler :=  http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
    serveMux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))
    serveMux.Handle("GET /api/healthz", http.HandlerFunc(handlerHeatlthz))
    serveMux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
    serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)

    server := http.Server{ Handler: serveMux, Addr: ":8080" }
    err := server.ListenAndServe()
    if err != nil {
        fmt.Printf("Error: %v", err) 
    }
}
