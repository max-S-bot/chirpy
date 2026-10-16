package main

import "net/http"

func main() {
    sm := http.NewServeMux()
    s := http.Server{
        Addr: ":8080",
        Handler: sm,
    }
    sm.Handle("/", http.FileServer(http.Dir(".")))
    sm.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Add("Content-Type", "text/plain; charset=utf-8") 
        w.WriteHeader(200)
        w.Write([]byte("OK"))
    })
    s.ListenAndServe()
}
