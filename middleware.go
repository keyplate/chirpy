package main

import (
    "net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        cfg.reqCount.Add(1)
        next.ServeHTTP(w, req)
    })
}
