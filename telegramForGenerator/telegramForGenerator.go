package telegramForGenerator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"wuzzapcom/bodiary/diaryGenerator"
	"wuzzapcom/bodiary/helpers"

	"github.com/Syfaro/telegram-bot-api"
)

const pathToUserDirectories = "" //"/Users/wuzzapcom/test/"

//TelegramForGenerator is main object
type TelegramForGenerator struct {
	bot     *tgbotapi.BotAPI
	updates <-chan tgbotapi.Update
	queue   []Queue
	users   map[int64]int
	mongo   *Mongo
}

//ConnectToTelegram Creates a connection to telegram, returns telegram object and channel with messages
func ConnectToTelegram() *TelegramForGenerator {
	log.Println("Start ConnectToTelegram")

	bot, err := tgbotapi.NewBotAPI(loadAuthToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	user := tgbotapi.NewUpdate(0)
	user.Timeout = 60

	updates, err := bot.GetUpdatesChan(user)
	if err != nil {
		log.Panic(err)
	}

	telegram := new(TelegramForGenerator)

	telegram.bot = bot
	telegram.updates = updates
	telegram.queue = make([]Queue, helpers.NumberOfThreads)
	telegram.users = make(map[int64]int)
	telegram.mongo = new(Mongo)
	telegram.mongo.init(telegram)
	for i := 0; i < helpers.NumberOfThreads; i++ {
		telegram.queue[i] = Queue{channel: make(chan tgbotapi.Update), telegram: telegram, users: make(map[int64]*helpers.UserValues)}
		go telegram.queue[i].workWithClient()
	}

	log.Println("End ConnectToTelegram")

	return telegram

}

func loadAuthToken() string {

	val, _ := ioutil.ReadFile(helpers.PathToAuthToken)

	str := string(val)

	return str[:len(str)-1]

}

func (telegram *TelegramForGenerator) loadBalancer(update tgbotapi.Update) {

	minNum := -1
	minVal := 0

	for i, queue := range telegram.queue {

		if minNum == -1 {
			minNum = i
			minVal = len(queue.channel)
		}

		length := len(queue.channel)
		if minVal > length {
			minVal = length
			minNum = i
		}

	}

	fmt.Printf("Message uploaded to clannel[%d]\n", minNum)

	telegram.queue[minNum].channel <- update

}

func (telegram *TelegramForGenerator) updateUserState(id int64, newState int) {

	telegram.users[id] = newState

}

func (telegram *TelegramForGenerator) getUserState(id int64) int {

	return telegram.users[id]

}

//WorkWithClient Main body which gets messages and handles them
func (telegram *TelegramForGenerator) WorkWithClient() {
	log.Println("Start WorkWithClient")

	for update := range telegram.updates {

		if update.Message == nil {
			continue
		}

		_, err := telegram.users[update.Message.Chat.ID]

		if err == false {
			telegram.users[update.Message.Chat.ID] = 0
		}

		telegram.loadBalancer(update)

	}

	log.Println("End WorkWithClient")

}

func (telegram *TelegramForGenerator) sendHTMLFileToUser(username string, id int64) {

	log.Println("Start sendHTMLFileToUser")

	userValues := telegram.mongo.getUserValuesByID(id)
	fmt.Println(userValues)
	diaryGenerator.GenerateDiary(userValues)
	telegram.bot.Send(tgbotapi.NewDocumentUpload(id, userValues.Name+".html"))
	os.Remove(userValues.Name + ".html")

	log.Println("End sendHTMLFileToUser")

}

func (telegram *TelegramForGenerator) openFile(username string, exp string) *os.File {

	log.Println("Start openFile")

	file, err := os.OpenFile(pathToUserDirectories+username+exp, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		fmt.Println(err)
	}

	log.Println("End openFile")

	return file

}

func (telegram *TelegramForGenerator) sendQueryToUser(id int64, message string) {
	log.Println("Start getOneQueryFromUser")

	telegram.bot.Send(tgbotapi.NewMessage(id, message))

	log.Println("End getOneQueryFromUser")

}

func (telegram *TelegramForGenerator) sendMessageUserDoesNotExit(userID int64) {
	log.Println("Start sendMessageUserDoesNotExit")

	telegram.bot.Send(tgbotapi.NewMessage(userID, "Пользователь не существует, воспользуйтесь командой /CreateNewUser"))

	log.Println("End sendMessageUserDoesNotExit")

}

func (telegram *TelegramForGenerator) checkUserRegistration(username string) bool {
	log.Println("Start checkUserRegistration")

	_, err := os.Stat(pathToUserDirectories + username + ".user")

	log.Println("End checkUserRegistration")

	return err != nil

}

//GetUserValues reades user file and returns special struct
func (telegram *TelegramForGenerator) GetUserValues(username string, id int64) *helpers.UserValues {

	log.Println("Start GetUserValues")

	//file := telegram.openFile()

	if telegram.checkUserRegistration(username) {

		telegram.sendMessageUserDoesNotExit(id)
		log.Println("End GetUserValues")
		return nil

	}

	data, _ := ioutil.ReadFile(pathToUserDirectories + username + ".user")

	start := 0
	userValues := new(helpers.UserValues)

	for i, val := range data {

		if val == byte('\n') {

			userValues.Name = string(data[start:i])
			start = i + 1
			break

		}

	}

	for i := start; i < len(data); i++ {

		if data[i] == byte('\n') {

			userValues.Group = string(data[start:i])
			start = i + 1
			break

		}

	}

	for i := start; i < len(data); i++ {

		if data[i] == byte('\n') {

			userValues.StartPulse, _ = strconv.Atoi(string(data[start:i]))
			start = i + 1
			break

		}

	}

	for i := start; i < len(data); i++ {

		if data[i] == byte('\n') {

			userValues.EndPulse, _ = strconv.Atoi(string(data[start:i]))
			break

		}

	}

	log.Println("End GetUserValues")

	return userValues

}
