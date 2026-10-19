package main

import (
    "net/http"
    "sync/atomic"
    "fmt"
    "encoding/json"
    "strings"
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
    sm.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, req *http.Request) {
        data := struct {Body string} {}
        json.NewDecoder(req.Body).Decode(&data) 
        if len(data.Body) > 140 {
            w.WriteHeader(400)
            w.Write([]byte("{\"error\": \"Chirp is too long\"}"))
        } else {
            ws := strings.Split(data.Body, " ")
            for i, w := range ws {
                w = strings.ToLower(w)
                if w == "kerfuffle" || w == "sharbert" || w == "fornax" {
                    ws[i] = "****"
                }
            }
            w.WriteHeader(200)
            w.Write([]byte("{\"cleaned_body\": \"" + strings.Join(ws, " ") + "\"}"))
        }
    })
    s.ListenAndServe()
}
