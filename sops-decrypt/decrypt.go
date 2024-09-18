package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
)

func manifestRequestHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/api/v1/getparams.execute" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        fmt.Fprintf(w, r.URL.Path)
        return
    }

    //if r.Method != "GET" {
    //    http.Error(w, r.Method, http.StatusNotFound)
    //    fmt.Fprintf(w, r.Method)
    //    return
    //}

    switch r.Method {
        case "POST":
            reqBody, err := ioutil.ReadAll(r.Body)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("%s\n", reqBody)
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