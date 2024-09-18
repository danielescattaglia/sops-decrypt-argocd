package main

import (
    "fmt"
    "log"
    "net/http"
)

func manifestRequestHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/api/v1/getparams.execute" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        fmt.Fprintf(w, r.URL.Path)
        return
    }

    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        fmt.Fprintf(w, r.Method)
        return
    }


    fmt.Fprintf(w, "Hello!")
}

func main() {

    http.HandleFunc("/api/v1/getparams.execute", manifestRequestHandler)

    fmt.Println("Hello, World!")
     if err := http.ListenAndServe(":4355", nil); err != nil {
        log.Fatal(err)
     }
}