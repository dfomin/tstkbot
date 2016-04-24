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
	"gopkg.in/mgo.v2/bson"
)

const (
	pigFileID  = "BQADAgAD6AAD9HsZAAF6rDKYKVsEPwI"
	dogeFileID = "BQADAgAD3gAD9HsZAAFphGBFqImfGAI"

	chickenNoFileID       = "BQADAgADswIAAkKvaQABArcCG5J-M4IC"
	chickenThinkingFileID = "BQADAgADvwIAAkKvaQABKt6_X0LBVfYC"
	chickenThumbUpFileID  = "BQADAgADnQIAAkKvaQABUb3ik6MhZwcC"
	chickenFacepalmFileID = "BQADAgADqwIAAkKvaQABEeHC3ECjvqwC"
	chickenWaitingFileID  = "BQADAgADsQIAAkKvaQAB72oOFFT5ryoC"
	chickenWhatFileID     = "BQADAgADvQIAAkKvaQABeZwGMzfkLroC"
	chickenWritingFileID  = "BQADAgADMQsAAkKvaQABwFPldEcMt14C"
	chickenDeadFileID     = "BQADAgADPQsAAkKvaQABih96aCmG-gQC"

	penguinDunnoFileID   = "BQADAQADyCIAAtpxZge0ITVcWNv_vwI"
	penguinLookOutFileID = "BQADAQADvCIAAtpxZgf5jpah4VvMqQI"

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

// From represents user who sended message
type From struct {
	ID int `json:"id"`
}

// Entity represents telegram message entity
type Entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

// Message represents telegram message info
type Message struct {
	From     From     `json:"from"`
	Chat     Chat     `json:"chat"`
	Text     string   `json:"text"`
	Entities []Entity `json:"entities,omitempty"`
}

// Object represents telegram message object
type Object struct {
	Message Message `json:"message"`
}

// JudgePhrase represents phrases for judging
type JudgePhrase struct {
	Phrase string `json:"phrase"`
}

// JudgePhraseCandidate represents candidate phrases for judging
type JudgePhraseCandidate struct {
	Phrase string `json:"phrase"`
	Users  []int  `json:"users"`
}

var mgoSession *mgo.Session

// InitDatabase represents database initialization
func InitDatabase(databaseName string) {
	var err error
	mgoSession, err = mgo.Dial(MongoDBHost)
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
		sendMessage(chat.ID, "тарахчу только в королях")
		sendSticker(chat.ID, chickenNoFileID)
		return
	}

	// Check command
	commandType, text := checkCommand(&object)
	fmt.Println(commandType)
	if commandType != "" {
		processCommand(commandType, text, &object)
	} else {
		processMessage(&object)
	}
}

func checkCommand(object *Object) (string, string) {
	for _, entity := range object.Message.Entities {
		if entity.Type == "bot_command" {
			command := object.Message.Text[entity.Offset : entity.Offset+entity.Length]
			text := ""
			if len(object.Message.Text) > entity.Offset+entity.Length+1 {
				text = object.Message.Text[entity.Offset+entity.Length+1:]
			}

			return command, text
		}
	}

	return "", object.Message.Text
}

func processCommand(command string, text string, object *Object) {
	if command == "/punto" || command == "/punto@TstkBot" {
		processPuntoCommand(object)
	} else if command == "/judge" || command == "/judge@TstkBot" {
		processJudgeCommand(object.Message.Chat.ID, text)
	} else if command == "/judgeadd" || command == "/judgeadd@TstkBot" {
		chatID := object.Message.Chat.ID
		phrase := text
		userID := object.Message.From.ID
		processJudgeAddCommand(chatID, phrase, userID)
	} else if command == "/judgeremove" || command == "/judgeremove@TstkBot" {
		processJudgeRemoveCommand()
	} else if command == "/judgelist" || command == "/judgelist@TstkBot" {
		processJudgeListCommand(object.Message.Chat.ID)
	}
}

func processPuntoCommand(object *Object) {
	count := rand.Intn(5) + 1
	for i := 0; i < count; i++ {
		sendSticker(object.Message.Chat.ID, pigFileID)
	}
}

func processJudgeCommand(chatID int, text string) {
	elements := strings.Split(text, " ")
	names := make([]string, len(elements))
	count := 0
	for _, element := range elements {
		name := strings.TrimSpace(element)
		if name != "" {
			names[count] = name
			count++
		}
	}

	names = names[:count]

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

	result := ""
	for _, name := range names {
		phrase := phrases[rand.Intn(len(phrases))].Phrase
		index := strings.Index(phrase, "#")

		if index == -1 {
			sendMessage(chatID, "что-то у фомы сломалось 😬😬😬")
			return
		}

		prefix := ""
		suffix := ""

		if index > 0 {
			prefix = phrase[:index-1]
		}

		if index < len(phrase)-1 {
			suffix = phrase[index+1:]
		}

		result += prefix + name + suffix + "\n"
	}

	sendMessage(chatID, result)
}

func processJudgeAddCommand(chatID int, phrase string, userID int) {
	sessionCopy := mgoSession.Copy()
	defer sessionCopy.Close()

	// Check that already added
	var phrases []JudgePhrase
	database := sessionCopy.DB(databaseName)
	phrasesCollection := database.C("judgePhrases")
	err := phrasesCollection.Find(bson.M{"phrase": phrase}).All(&phrases)
	if err != nil {
		sendMessage(chatID, "что-то у фомы сломалось 😬😬😬")
		return
	}

	if len(phrases) > 0 {
		sendSticker(chatID, chickenWhatFileID)
		return
	}

	var candidates []JudgePhraseCandidate
	candidatesCollection := database.C("judgePhrasesCandidates")
	err = candidatesCollection.Find(bson.M{"phrase": phrase}).All(&candidates)
	if err != nil {
		sendMessage(chatID, "что-то у фомы сломалось 😬😬😬")
		return
	}

	if len(candidates) > 1 {
		sendMessage(chatID, "что-то у фомы сломалось 😬😬😬")
		return
	}

	var candidate JudgePhraseCandidate
	if len(candidates) == 0 {
		var newCandidate JudgePhraseCandidate
		newCandidate.Phrase = phrase
		newCandidate.Users = make([]int, 3)
		candidate = newCandidate
	} else {
		candidate = candidates[0]
	}

	// Add new user
	for i := 0; i < 3; i++ {
		if candidate.Users[i] == 0 {
			candidate.Users[i] = userID
			break
		} else if candidate.Users[i] == userID {
			break
		}
	}

	// Count how much users
	count := 3
	for i := 0; i < 3; i++ {
		if candidate.Users[i] == 0 {
			count = i
			break
		}
	}

	// TODO: check errors
	switch count {
	case 1:
		_, err = candidatesCollection.Upsert(bson.M{"phrase": phrase}, candidate)
		sendSticker(chatID, chickenNoFileID)
	case 2:
		_, err = candidatesCollection.Upsert(bson.M{"phrase": phrase}, candidate)
		sendSticker(chatID, chickenThinkingFileID)
	case 3:
		err = candidatesCollection.Remove(bson.M{"phrase": phrase})
		var newPhrase JudgePhrase
		newPhrase.Phrase = phrase
		_, err = phrasesCollection.Upsert(bson.M{"phrase": phrase}, newPhrase)
		sendSticker(chatID, chickenWritingFileID)
	}
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
