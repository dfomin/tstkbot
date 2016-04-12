package main

import (
    "fmt"
    "net/http"
    //"net/http/httputil"
    "net/url"
    "encoding/json"
    "strconv"
    "strings"
)

const (
    pigFileId = "BQADAgAD6AAD9HsZAAF6rDKYKVsEPwI"
)

type Chat struct {
    Id int `json:"id"`
}

type Message struct {
    Chat Chat `json:"chat"`
    Text string `json:"text"`
}

type Object struct {
    Message Message `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    //data, _ := httputil.DumpRequest(r, true)
    //fmt.Printf("%s\n\n", data)
    var object Object
    err := json.NewDecoder(r.Body).Decode(&object)
    if err != nil {
        fmt.Println(err)
    }

    index := strings.Index(object.Message.Text, "/punto")
    if index != -1 {
        sendSticker(object.Message.Chat.Id, pigFileId)
    } else {
        sendMessage(object.Message.Chat.Id, object.Message.Text)
    }
}

func sendMessage(id int, text string) {
    apiUrl := `https://api.telegram.org/bot120816766:AAHuy66RPZLVt3JwBWPwGh2Ndxt_KwAXYlE/sendMessage`
    data := url.Values{}
    data.Add("chat_id", strconv.Itoa(id))
    data.Add("text", text)

    resp, err := http.Get(apiUrl + "?" + data.Encode())
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
}

func sendSticker(id int, fileId string) {
    apiUrl := `https://api.telegram.org/bot120816766:AAHuy66RPZLVt3JwBWPwGh2Ndxt_KwAXYlE/sendSticker`
    data := url.Values{}
    data.Add("chat_id", strconv.Itoa(id))
    data.Add("sticker", fileId)

    //fmt.Println(data.Encode())
    resp, err := http.Get(apiUrl + "?" + data.Encode())
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
}

func main() {
    http.HandleFunc("/tstkbot", handler)
    err := http.ListenAndServeTLS(":8443", "fullchain.pem", "privkey.pem", nil)
    if err != nil {
        fmt.Println("ListenAndServe: ", err)
    }
}

