package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "strings"

    "github.com/getsops/sops/v3/decrypt"
    "github.com/getsops/sops/v3/cmd/sops/formats"
)

func decryptFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading %q: %w", file, err)
	}

	format := formats.FormatForPath(file)
	data, err := decrypt.DataWithFormat(b, format)
	if err != nil {
		return nil, fmt.Errorf("trouble decrypting file: %w", err)
	}
	return data, nil
}

func healthzRequestHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "GET":
            if r.URL.Path == "/healthz" {
                w.Write ([]byte(`OK`))
            } else {
                fmt.Println(r.URL.Path)
            }
        default:
            fmt.Println(r.URL.Path)
    }
}

func manifestRequestHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
        case "POST":
            if r.URL.Path == "/api/v1/getparams.execute" {
                reqBody, err := ioutil.ReadAll(r.Body)
                if err != nil {
                    log.Fatal(err)
                }
                //fmt.Fprintf(w, "{ \"output\": { \"valuesObject\": [ { \"keyrenewperiod\": \"10\", } ] } }")
                w.Header().Set("Content-Type", "application/json")

                jsonData := createBody()
                w.Write (jsonData)

                fmt.Println(string(reqBody))
            } else {
                http.Error(w, "404 not found.", http.StatusNotFound)
                fmt.Fprintf(w, r.URL.Path)
                return
            }
        default:
            fmt.Println(r.URL.Path)
    }
}

//
func createBody() ([]byte) {
    data, err2 := decryptFile("./test-values.yaml")

    jsonData := []byte(`{
    "output": {
        "parameters": [
            {
                "valuesobject": "` + strings.Replace(string(data), "\"", "\\\"", -1) + `"
            }
        ]
    }
}`)

    fmt.Println(string(jsonData))
    fmt.Println(err2)

    return jsonData
}

func main() {

    //fmt.Println(string(createBody()))

    http.HandleFunc("/api/v1/getparams.execute", manifestRequestHandler)
    http.HandleFunc("/healthz", healthzRequestHandler)

    fmt.Println("Hello, World!")
    if err := http.ListenAndServe(":4355", nil); err != nil {
        log.Fatal(err)
    }
}
