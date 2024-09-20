package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"

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
            //fmt.Fprintf(w, "{ \"output\": { \"valuesObject\": [ { \"keyrenewperiod\": \"10\", } ] } }")
            w.Header().Set("Content-Type", "application/json")

            jsonData := createBody()
            w.Write (jsonData)

            fmt.Println(string(reqBody))
    }
}

func createBody() ([]byte) {
    data, err2 := decryptFile("./values.yaml")

    jsonData := []byte(`{
    "output": {
        "parameters": [
            {
                "valuesobject": "` + string(data) + `"
            }
        ]
    }
}`)

    fmt.Println(string(jsonData))
    fmt.Println(err2)

    return jsonData
}

func main() {

    //fmt.Println(createBody())

    http.HandleFunc("/api/v1/getparams.execute", manifestRequestHandler)

    fmt.Println("Hello, World!")
    if err := http.ListenAndServe(":4355", nil); err != nil {
        log.Fatal(err)
    }
}
