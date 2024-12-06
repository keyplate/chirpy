package main

import (
    "fmt"
    "net/http"
)

func main() {
    serveMux := http.NewServeMux()
    serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
    serveMux.Handle("/healthz", http.HandlerFunc(handleHeatlthz))
    server := http.Server{ Handler: serveMux, Addr: ":8080" }
    err := server.ListenAndServe()
    if err != nil {
        fmt.Printf("Error: %v", err) 
    }
}
