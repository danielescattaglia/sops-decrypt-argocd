package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "strings"
    "encoding/json"

    "github.com/getsops/sops/v3/decrypt"
    "github.com/getsops/sops/v3/cmd/sops/formats"
)

type argocdAppParams struct {
	Input Input `json:"input"`
}
type Parameters struct {
	ObjectType    string `json:"objectType"`
	EncryptedFile string `json:"encryptedFile"`
}
type Input struct {
	Parameters Parameters `json:"parameters"`
}

func decryptFile(file string) ([]byte, error) {
    b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading %q: %w", file, err)
	}

	format := formats.FormatForPath(file)

	data, err := decryptBytes(b, format)

	return data, err
}

func decryptContent(content string) ([]byte, error) {
    format := formats.FormatFromString(content)

	data, err := decryptBytes([]byte(content), format)

	return data, err
}

func decryptBytes(b []byte, format formats.Format) ([]byte, error) {
	data, err := decrypt.DataWithFormat(b, format)
	if err != nil {
		return nil, fmt.Errorf("trouble decrypting file: %w", err)
	}
	return data, nil
}

func createBody(encryptedContent string) ([]byte) {
    data, err2 := decryptContent(encryptedContent)

    if err2 != nil {
        jsonData := []byte(`{
        "output": {
            "parameters": [
                {
                    "valuesobject": "` + strings.Replace(string(data), "\"", "\\\"", -1) + `"
                }
            ]
        }
    }`)

        return jsonData
    } else {
        fmt.Println(err2)
    }

    return nil
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

                var jsonBody argocdAppParams
                errUnmarshal := json.Unmarshal(reqBody, &jsonBody)
                if errUnmarshal != nil {
                    log.Fatalf("Unable to marshal JSON due to %s", errUnmarshal)
                }

                jsonData := createBody(jsonBody.Input.Parameters.EncryptedFile)
                //jsonData :=  []byte(`{ "output": { "valuesObject": { \"keyrenewperiod\": \"10\", } } }`)
fmt.Println(string(jsonData))
                w.Header().Set("Content-Type", "application/json")
                w.Write (jsonData)

            } else {
                http.Error(w, "404 not found.", http.StatusNotFound)
                fmt.Fprintf(w, r.URL.Path)
                return
            }
        default:
            fmt.Println(r.URL.Path)
    }
}

func test (w http.ResponseWriter, r *http.Request) {
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
    }

    //var jsonBodyMap map[string]interface{}
    var reqInput argocdAppParams

    err2 := json.Unmarshal(reqBody, &reqInput)

    fmt.Println(err2)
    fmt.Println(reqInput.Input.Parameters.ObjectType)

    //for k, v := range jsonBodyMap["parameters"] {
    //    fmt.Printf("key[%s] value[%s]\n", k, v)
    //}
}

func main() {

    // in questo momento mi faccio passare da decriptare come stringa tra i parametri di input
    // in flux come funziona? Va a leggere il file in un repository?

    //fmt.Println(string(createBody()))

    http.HandleFunc("/api/v1/getparams.execute", manifestRequestHandler)
    http.HandleFunc("/healthz", healthzRequestHandler)
    http.HandleFunc("/test", test)

    fmt.Println("Hello, World!")
    if err := http.ListenAndServe(":4355", nil); err != nil {
        log.Fatal(err)
    }
}
