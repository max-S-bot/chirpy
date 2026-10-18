package main

import (
    "net/http"
    "sync/atomic"
    "fmt"
)

func main() {
    sm := http.NewServeMux()
    s := http.Server{
        Addr: ":8000",
        Handler: sm,
    }
    hits := &atomic.Int32{}
    handler := http.StripPrefix("/app", http.FileServer(http.Dir("./app")))
    sm.HandleFunc("/app/", func (w http.ResponseWriter, req *http.Request) {
        hits.Add(1)
        handler.ServeHTTP(w, req)
    })
    sm.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Add("Content-Type", "text/plain; charset=utf-8") 
        w.WriteHeader(200)
        w.Write([]byte("OK"))
    })
    sm.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Add("Content-Type", "text/html")
        w.Write(fmt.Appendf([]byte{}, `<html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>
        </html>
        `, hits.Load()))
    })
    sm.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, req *http.Request) {
        hits.Store(0)
    })
    s.ListenAndServe()
}
