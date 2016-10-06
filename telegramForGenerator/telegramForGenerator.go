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
	bot      *tgbotapi.BotAPI
	updates  <-chan tgbotapi.Update
	queue    []Queue
	users    map[int64]int
	reminder helpers.RemiderDates
}

//ConnectToTelegram Creates a connection to telegram, returns telegram object and channel with messages
func ConnectToTelegram() *TelegramForGenerator {
	log.Println("Start ConnectToTelegram")

	bot, err := tgbotapi.NewBotAPI("263647981:AAGCsCIqCVi3c089IJ47oLWc26Ix9SvbuPE")
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
	telegram.loadReminder()
	telegram.queue = make([]Queue, helpers.NumberOfThreads)
	for i := 0; i < helpers.NumberOfThreads; i++ {
		telegram.queue[i] = Queue{channel: make(chan tgbotapi.Update), telegram: telegram}
		go telegram.queue[i].workWithClient()
	}

	log.Println("End ConnectToTelegram")

	return telegram

}

func (telegram *TelegramForGenerator) loadReminder() {

	log.Println("Start loadReminder")

	_, err := os.Stat(pathToUserDirectories + "reminderDates.serialize")

	if err == nil {

		data, _ := ioutil.ReadFile(pathToUserDirectories + "reminderDates.serialize")
		telegram.reminder = helpers.FromGOB64(string(data))

	} else {

		telegram.reminder = make(helpers.RemiderDates) //*(new(helpers.RemiderDates))
		for i := 0; i < 7; i++ {
			telegram.reminder[i] = make([]helpers.Pair, 0)
		}
		os.Create(pathToUserDirectories + "reminderDates.serialize")

	}

	log.Println("End loadReminder")

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
		if err == true {
			telegram.users[update.Message.Chat.ID] = 0
		}

		telegram.loadBalancer(update)

	}

	log.Println("End WorkWithClient")

}

func (telegram *TelegramForGenerator) sendHTMLFileToUser(username string, id int64) {

	log.Println("Start sendHTMLFileToUser")
	userValues := telegram.GetUserValues(username, id)
	if userValues != nil {
		diaryGenerator.GenerateDiary(userValues)
		telegram.bot.Send(tgbotapi.NewDocumentUpload(id, pathToUserDirectories+username+".html"))
	}

	log.Println("End sendHTMLFileToUser")

}

//createFile tested
func (telegram *TelegramForGenerator) createFile(username string) *os.File {
	log.Println("Start createFile")

	file, _ := os.Create(pathToUserDirectories + username + ".user")

	log.Println("End createFile")

	return file

}

func (telegram *TelegramForGenerator) openFile(username string, exp string) *os.File {

	log.Println("Start openFile")

	file, err := os.Open(pathToUserDirectories + username + exp)

	if err != nil {
		log.Fatal(err)
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

func (telegram *TelegramForGenerator) saveReminder() {

	log.Println("Start saveReminder")

	reminderFile, err := os.OpenFile(pathToUserDirectories+"reminderDates.serialize", os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err.Error())
	}
	reminderFile.WriteString(helpers.ToGOB64(telegram.reminder))
	reminderFile.Close()

	//.Println(helpers.ToGOB64(telegram.reminder))

	log.Println("End saveReminder")

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

	userValues.UserName = username

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

//===================
//REMINDER-NOT FINISHED

// func (telegram *TelegramForGenerator) sendReminderInfinityLoop() {
// 	log.Println("Start sendReminderInfinityLoop")
//
// 	for {
// 		//telegram.sendRemind()
// 		telegram.sendRemindToUser(helpers.Pair{ID: telegram.id, UserName: telegram.username, RemindThisWeek: true})
// 		break
//
// 	}
//
// 	log.Println("End sendReminderInfinityLoop")
//
// }
//
// func (telegram *TelegramForGenerator) sendRemind() {
// 	log.Println("Start sendRemind")
//
// 	currentDayOfWeek := (int(time.Now().Weekday()) + 6) % 7
//
// 	if len(telegram.reminder[currentDayOfWeek]) != 0 && time.Now().Hour() > 17 {
//
// 		for i, val := range telegram.reminder[currentDayOfWeek] {
//
// 			if val.RemindThisWeek {
//
// 				telegram.sendRemindToUser(val)
// 				telegram.reminder[currentDayOfWeek][i].RemindThisWeek = false
//
// 			} else {
// 				telegram.reminder[currentDayOfWeek][i].RemindThisWeek = true
// 			}
//
// 		}
//
// 	}
//
// 	log.Println("End sendRemind")
//
// }
//
// func (telegram *TelegramForGenerator) sendRemindToUser(user helpers.Pair) {
// 	log.Println("Start sendRemindToUser")
//
// 	text := fmt.Sprintf("Привет, %s! Пришло время сдать дневник самоподготовки.", user.UserName)
// 	message := tgbotapi.NewMessage(user.ID, text)
//
// 	telegram.bot.Send(message)
//
// 	telegram.sendHTMLFileToUser(user.ID, user.UserName)
//
// 	log.Println("End sendRemindToUser")
//
// }
