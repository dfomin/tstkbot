package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"encoding/json"
	"math/rand"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	pigFileID             = "BQADAgAD6AAD9HsZAAF6rDKYKVsEPwI"
	dogeFileID            = "BQADAgAD3gAD9HsZAAFphGBFqImfGAI"
	chickenNoFileID       = "BQADAgADswIAAkKvaQABArcCG5J-M4IC"
	chickenThinkingFileID = "BQADAgADvwIAAkKvaQABKt6_X0LBVfYC"
	chickenThumbUpFileID  = "BQADAgADnQIAAkKvaQABUb3ik6MhZwcC"
	penguinDunnoFileID    = "BQADAQADyCIAAtpxZge0ITVcWNv_vwI"
	penguinLookOutFileID  = "BQADAQADvCIAAtpxZgf5jpah4VvMqQI"

	apiURL = `https://api.telegram.org/bot120816766:AAHuy66RPZLVt3JwBWPwGh2Ndxt_KwAXYlE/`

	// MongoDBHost represents mongo db host and port
	MongoDBHost  = "127.0.0.1:27017"
	databaseName = "tstkbot"

	tstkChatID = -14369410
)

// Chat represents telegram chat info
type Chat struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

// Entity represents telegram message entity
type Entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

// Message represents telegram message info
type Message struct {
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
	Entities []Entity `json:"entities,omitempty"`
}

// Object represents telegram message object
type Object struct {
	Message Message `json:"message"`
}

type JudgePhrase struct {
	Phrase string `json:"phrase"`
}

var mgoSession *mgo.Session

// InitDatabase represents database initialization
func InitDatabase(databaseName string) {
	mgoSession, err := mgo.Dial(MongoDBHost)
	if err != nil {
		fmt.Println("Failed to connect to database")
		log.Fatal(err)
	}

	mgoSession.SetMode(mgo.Monotonic, true)
}

// Process user message
func gotMessage(w http.ResponseWriter, r *http.Request) {
	data, _ := httputil.DumpRequest(r, true)
	fmt.Printf("%s\n\n", data)

	// Parse telegram message
	var object Object
	err := json.NewDecoder(r.Body).Decode(&object)
	if err != nil {
		fmt.Println(err)
	}

	// Check chat, only tstk chat is supported
	chat := object.Message.Chat
	if chat.ID != tstkChatID && chat.Type == "group" {
		sendMessage(chat.ID, "Тарахчу только в королях")
		sendSticker(chat.ID, chickenNoFileID)
		return
	}

	// Check command
	commandType := checkCommand(&object)
	fmt.Println(commandType)
	if commandType != "" {
		processCommand(commandType, &object)
	} else {
		processMessage(&object)
	}
}

func checkCommand(object *Object) string {
	for _, entity := range object.Message.Entities {
		if entity.Type == "bot_command" {
			return object.Message.Text[entity.Offset : entity.Offset+entity.Length]
		}
	}

	return ""
}

func processCommand(command string, object *Object) {
	if command == "/punto" || command == "/punto@TstkBot" {
		processPuntoCommand(object)
	} else if command == "/judge" || command == "/judge@TstkBot" {
		// TODO: fix
		names := strings.Split(object.Message.Text, " ")
		if len(names) > 1 {
			processJudgeCommand(object.Message.Chat.ID, names[1:])
		} else {
			sendMessage(object.Message.Chat.ID, "бесишь")
		}
	} else if command == "/judgeAdd" || command == "/judgeAdd@TstkBot" {
		processJudgeAddCommand()
	} else if command == "/judgeRemove" || command == "/judgeRemove@TstkBot" {
		processJudgeRemoveCommand()
	} else if command == "/judgeList" || command == "/judgeList@TstkBot" {
		processJudgeListCommand(object.Message.Chat.ID)
	}
}

func processPuntoCommand(object *Object) {
	count := rand.Intn(5) + 1
	for i := 0; i < count; i++ {
		sendSticker(object.Message.Chat.ID, pigFileID)
	}
}

func processJudgeCommand(id int, names []string) {
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

func processJudgeAddCommand() {

}

func processJudgeRemoveCommand() {

}

func processJudgeListCommand(chatID int) {
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()

	var phrases []JudgePhrase
	database := sessionCopy.DB(databaseName)
	phrasesCollection := database.C("judgePhrases")
	err := phrasesCollection.Find(nil).All(&phrases)
	if err != nil {
		sendMessage(chatID, "что-то у фомы сломалось 😬😬😬")
		return
	}

	if len(phrases) == 0 {
		sendSticker(chatID, penguinDunnoFileID)
		return
	}

	answer := ""
	for _, phrase := range phrases {
		answer += phrase.Phrase + "\n"
	}

	sendMessage(chatID, answer)
}

func processMessage(object *Object) {
	text := object.Message.Text
	if text == "" {
		// TODO: empty message
	} else if string(text[len(text)-1]) == "?" {
		processQuestionMessage(object)
	} else {
		processStatementMessage(object)
	}

	//sendMessage(object.Message.Chat.ID, selectAnswer())
}

func processQuestionMessage(object *Object) {
	sendMessage(object.Message.Chat.ID, "не знаю, фома не накодил")
}

func processStatementMessage(object *Object) {
	sendMessage(object.Message.Chat.ID, "фома не накодил")
}

func dogeSender(id int) {
	delay := rand.Intn(60*4) + 60*8
	time.Sleep(time.Duration(delay) * time.Minute)
	sendSticker(id, dogeFileID)
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
	answerURL := apiURL + "sendMessage"
	data := url.Values{}
	data.Add("chat_id", strconv.Itoa(id))
	data.Add("text", text)

	resp, err := http.Get(answerURL + "?" + data.Encode())
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func sendSticker(id int, fileID string) {
	answerURL := apiURL + "sendSticker"
	data := url.Values{}
	data.Add("chat_id", strconv.Itoa(id))
	data.Add("sticker", fileID)

	resp, err := http.Get(answerURL + "?" + data.Encode())
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
}

func main() {
	InitDatabase(databaseName)

	http.HandleFunc("/tstkbot", gotMessage)
	err := http.ListenAndServeTLS(":8443", "fullchain.pem", "privkey.pem", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
