package main

import (
    "fmt"
    "net/http"
    "net/http/httputil"
    "io/ioutil"
    "net/url"
)

func handler(w http.ResponseWriter, r *http.Request) {
    data, _ := httputil.DumpRequest(r, true)
    fmt.Printf("%s\n\n", data)

    sendMessage()
}

func sendMessage() {
    apiUrl := `https://api.telegram.org/bot120816766:AAHuy66RPZLVt3JwBWPwGh2Ndxt_KwAXYlE/sendMessage`
    data := url.Values{}
    data.Add("chat_id", "45227519")
    data.Add("text", "что в итоге?")

    resp, err := http.Get(apiUrl + "?" + data.Encode())
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}

func main() {
    http.HandleFunc("/tstkbot", handler)
    err := http.ListenAndServeTLS(":8443", "fullchain.pem", "privkey.pem", nil)
    if err != nil {
        fmt.Println("ListenAndServe: ", err)
    }
}

