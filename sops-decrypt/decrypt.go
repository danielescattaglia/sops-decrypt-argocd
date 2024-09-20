package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "strings"
    "encoding/json"
    //"bytes"

    "github.com/getsops/sops/v3/decrypt"
    "github.com/getsops/sops/v3/cmd/sops/formats"
)

type argocdAppParams struct {
	Input Input `json:"input"`
}
type Parameters struct {
	ObjectType    string `json:"objectType"`
	EncryptedFile string `json:"encryptedFile"`
	EncryptedFileType string `json:"encryptedFileType"`
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
    format := formats.FormatFromString("yaml")

	data, err := decryptBytes([]byte(content), format)

    if err != nil {
        fmt.Println(err)
        return nil, err
    }

	return data, nil
}

func decryptBytes(b []byte, f formats.Format) ([]byte, error) {
	data, err := decrypt.DataWithFormat(b, f)

	if err != nil {
		return nil, fmt.Errorf("trouble decrypting file: %w", err)
	}
	return data, nil
}

func jsonEscape(j string) string {
    var dataString string

	dataString = strings.Replace(j, `\`, `\\`, -1)
    dataString = strings.Replace(dataString, "\n", `\n`, -1)
    dataString = strings.Replace(dataString, "\b", `\b`, -1)
    dataString = strings.Replace(dataString, "\f", `\f`, -1)
    dataString = strings.Replace(dataString, "\r", `\r`, -1)
    dataString = strings.Replace(dataString, "\t", `\t`, -1)
    dataString = strings.Replace(dataString, "\"", "\\\"", -1)

    return dataString
}

func createBody(encryptedContent string, encryptedType string) ([]byte) {
    var data []byte = nil
    var err error = nil

    switch encryptedType {
        case "content":
            data, err = decryptContent(encryptedContent)
        case "file":
            data, err = decryptFile(encryptedContent)
    }

    if err != nil {
        fmt.Println(err)
        return nil
    }

    dataString := jsonEscape(string(data))

    jsonData := []byte(`{
       "output": {
           "parameters": [
               {
                   "valuesobject": "` + dataString + `"
               }
           ]
       }
   }`)

   return jsonData
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

func getAuthorizationToken (file string) string {
    b, err := ioutil.ReadFile(file)
	if err != nil {
	    fmt.Errorf("error getting authorization token %q: %w", file, err)
		return ""
	}

    return string(b)
}

func manifestRequestHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get ("Authorization") != "Bearer " + getAuthorizationToken("/var/run/argo/token") {
        fmt.Println("Token different, cannot proceed")
        return
    }

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
                    return
                }

                jsonData := createBody(jsonBody.Input.Parameters.EncryptedFile, jsonBody.Input.Parameters.EncryptedFileType)
                //jsonData :=  []byte(`{ "output": { "valuesObject": { \"keyrenewperiod\": \"10\", } } }`)

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

    b2, err := ioutil.ReadFile("./test-secret.yaml")
	if err != nil {
		fmt.Println (fmt.Errorf("error reading: %w", err))
	}

    fmt.Println(strings.Replace(string(b2), "\n", `\n`, -1))
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
