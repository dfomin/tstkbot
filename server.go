package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"net/http/httputil"
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	pigFileID  = "BQADAgAD6AAD9HsZAAF6rDKYKVsEPwI"
	dogeFileID = "BQADAgAD3gAD9HsZAAFphGBFqImfGAI"

	apiURL = `https://api.telegram.org/bot120816766:AAHuy66RPZLVt3JwBWPwGh2Ndxt_KwAXYlE/`

	// MongoDBHost represents mongo db host and port
	MongoDBHost = "127.0.0.1:27017"
)

// Chat represents telegram chat info
type Chat struct {
	ID int `json:"id"`
}

// Entities represents telegram message entities
type Entities struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

// Message represents telegram message info
type Message struct {
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
	Entities Entities `json:"entities,omitempty"`
}

// Object represents telegram message object
type Object struct {
	Message Message `json:"message"`
}

// Controller represents controller for database
type Controller struct {
	DatabaseName string
	Session      *mgo.Session
}

var mgoSession *mgo.Session

// InitDatabase represents database initialization
func InitDatabase(databaseName string) {
	mgoSession, err := mgo.Dial(MongoDBHost)
	if err != nil {
		fmt.Println("Failed to connect to database")
		log.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)
}

// Process user message
func gotMessage(w http.ResponseWriter, r *http.Request) {
	// Parse telegram message
	var object Object
	err := json.NewDecoder(r.Body).Decode(&object)
	if err != nil {
		fmt.Println(err)
	}

	// Check command
	commandType := checkCommand(object)
    fmt.Println(commandType)
	if commandType != "" {
		processCommand(commandType, object)
	} else {
		processMessage(object)
	}
}

func checkCommand(object *Object) string {
    if object.Message.Entities.Type == "" {
        return ""
    } else {
        type := object.Message.Entities.Type
        offset := object.Message.Entities.Offset
        length := object.Message.Entities.Length
        return type[offset:offset+length]
    }
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
	delay := rand.Intn(60*4) + 60*8
	time.Sleep(time.Duration(delay) * time.Minute)
	sendSticker(id, dogeFileId)
	go dogeSender(id)
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
		"чушь",
	}

	return answers[rand.Intn(len(answers))]
}

// Send commands

func sendMessage(id int, text string) {
	url := apiURL + "sendMessage"
	data := url.Values{}
	data.Add("chat_id", strconv.Itoa(id))
	data.Add("text", text)

	resp, err := http.Get(apiUrl + "?" + data.Encode())
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func sendSticker(id int, fileID string) {
	url := apiURL + "sendSticker"
	data := url.Values{}
	data.Add("chat_id", strconv.Itoa(id))
	data.Add("sticker", fileId)

	resp, err := http.Get(url + "?" + data.Encode())
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
}

func main() {
	dc := database.InitDatabase("tstkbot")

	http.HandleFunc("/tstkbot", gotMessage)
	err := http.ListenAndServeTLS(":8443", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
