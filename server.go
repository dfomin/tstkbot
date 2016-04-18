package main

import (
    "fmt"
    "net/http"
    //"net/http/httputil"
    "net/url"
    "encoding/json"
    "strconv"
    "strings"
    "time"
    "math/rand"
)

const (
    pigFileId = "BQADAgAD6AAD9HsZAAF6rDKYKVsEPwI"
    dogeFileId = "BQADAgAD3gAD9HsZAAFphGBFqImfGAI"
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

var dogeSubscription = false

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
        return
    }

    index = strings.Index(object.Message.Text, "/judge")
    if index != -1 {
        names := strings.Split(object.Message.Text, " ")
        if len(names) > 1 {
            judge(object.Message.Chat.Id, names[1:])
        } else {
            sendMessage(object.Message.Chat.Id, "бесишь")
        }
        return
    }

    index = strings.Index(object.Message.Text, "/doge")
    if index != -1 {
        if !dogeSubscription {
            dogeSubscription = true
            go dogeSender(object.Message.Chat.Id)
        }
        return
    }

    sendMessage(object.Message.Chat.Id, selectAnswer())
}

func judge(id int, names []string) {
    phrases := []string{
        "ноет",
        "по делу",
        "не по делу",
        "развернул шатер",
        "клоун",
        "без нытья",
        "сел в лужу",
        "кромсает",
        "уничтожил на молекулы",
        "перебор",
        "что-то ага",
        "ни в какие ворота",
        "самоуничтожился",
        "несет чушь",
        "жопка в тепле",
        "бесит",
        "байки травит",
        "подгорел",
        "бомбанул",
        "отскок",
        "устроил срач",
        "куда полез?",
    }
    var result string
    for _, name := range names {
        phrase := phrases[rand.Intn(len(phrases))]
        result += name + " " + phrase + ", "
    }

    sendMessage(id, result[:len(result)-2])
}

func dogeSender(id int) {
    delay := rand.Intn(60 * 4) + 60 * 8
    time.Sleep(time.Duration(delay) * time.Minute)
    sendSticker(id, dogeFileId)
    go dogeSender(id)
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

func selectAnswer() string {
    answers := []string{
        "и че?",
        "ну и?",
        "тебя не спросил",
        "отр",
        "о том и речь",
        "ноешь",
        "байки",
        "хм",
        "хер знает",
        "😬😬😬",
        "офк",
    }

    return answers[rand.Intn(len(answers))]
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

