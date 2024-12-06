package main

import (
	"net/http"
)

func handleHeatlthz(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte("OK"))
    w.WriteHeader(http.StatusOK)
} 
