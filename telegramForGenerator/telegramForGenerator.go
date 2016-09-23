package telegramForGenerator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
	"wuzzapcom/bodiary/diaryGenerator"
	"wuzzapcom/bodiary/helpers"

	"github.com/Syfaro/telegram-bot-api"
)

const pathToUserDirectories = "" //"/Users/wuzzapcom/test/"

//TelegramForGenerator is main object
type TelegramForGenerator struct {
	bot           *tgbotapi.BotAPI
	updates       <-chan tgbotapi.Update
	id            int64
	username      string
	reminder      helpers.RemiderDates
	multiThreader map[int]chan tgbotapi.Update
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

//WorkWithClient Main body which gets messages and handles them
func (telegram *TelegramForGenerator) WorkWithClient() {
	log.Println("Start WorkWithClient")

	for update := range telegram.updates {

		if update.Message == nil {
			continue
		}

		fmt.Println("Start cycle WorkWithClient")

		telegram.username = update.Message.Chat.UserName
		telegram.id = update.Message.Chat.ID

		if update.Message.Command() == "CreateNewUser" {

			telegram.getBasicInformationFromClient(update)
			diaryGenerator.GenerateDiary(telegram.GetUserValues())

		} else if update.Message.Command() == "GetDiary" {

			telegram.sendHTMLFileToUser(telegram.id, telegram.username)

		} else if update.Message.Command() == "help" {

			telegram.bot.Send(tgbotapi.NewMessage(telegram.id, helpers.Help))

		}

		fmt.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		telegram.bot.Send(msg)

		fmt.Println("End cycle WorkWithClient")

	}

	log.Println("End WorkWithClient")

}

func (telegram *TelegramForGenerator) sendHTMLFileToUser(id int64, username string) {

	log.Println("Start sendHTMLFileToUser")
	userValues := telegram.GetUserValues()
	if userValues != nil {
		diaryGenerator.GenerateDiary(telegram.GetUserValues())
		telegram.bot.Send(tgbotapi.NewDocumentUpload(id, pathToUserDirectories+username+".html"))
	}

	log.Println("End sendHTMLFileToUser")

}

//createFile tested
func (telegram *TelegramForGenerator) createFile() *os.File {
	log.Println("Start createFile")

	file, _ := os.Create(pathToUserDirectories + telegram.username + ".user")

	log.Println("End createFile")

	return file

}

func (telegram *TelegramForGenerator) openFile(exp string) *os.File {

	log.Println("Start openFile")

	file, err := os.Open(pathToUserDirectories + telegram.username + exp)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("End openFile")

	return file

}

//getNextMessage tested
func (telegram *TelegramForGenerator) getNextMessage() (tgbotapi.Update, error) {

	log.Println("Start getNextMessage")
	for update := range telegram.updates {

		if update.Message == nil {
			continue
		}

		log.Println("End getNextMessage")
		return update, nil

	}

	upd := new(tgbotapi.Update)

	log.Println("End getNextMessage")
	return *upd, errors.New("empty message")

}

func (telegram *TelegramForGenerator) sendMessageUserDoesNotExit() {
	log.Println("Start sendMessageUserDoesNotExit")

	telegram.bot.Send(tgbotapi.NewMessage(telegram.id, "Пользователь не существует, воспользуйтесь командой /CreateNewUser"))

	log.Println("End sendMessageUserDoesNotExit")

}

func (telegram *TelegramForGenerator) getOneQueryFromUser(update tgbotapi.Update, message string, file *os.File) {
	log.Println("Start getOneQueryFromUser")

	telegram.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
	update, _ = telegram.getNextMessage()
	file.WriteString(update.Message.Text + "\n")

	log.Println("End getOneQueryFromUser")

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

//getBasicInformationFromClient tested
func (telegram *TelegramForGenerator) getBasicInformationFromClient(update tgbotapi.Update) {

	log.Println("Start getBasicInformation")

	file := telegram.createFile()

	telegram.reminder[1] = append(telegram.reminder[1], helpers.Pair{ID: telegram.id, UserName: telegram.username, RemindThisWeek: true})

	telegram.getOneQueryFromUser(update, "Введите имя студента", file)

	telegram.getOneQueryFromUser(update, "Введите группу", file)

	telegram.getOneQueryFromUser(update, "Введите начальный пульс", file)

	telegram.getOneQueryFromUser(update, "Введите конечный пульс", file)

	telegram.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Готово!"))

	file.Close()

	telegram.saveReminder()

	log.Println("End getBasicInformation")

}

func (telegram *TelegramForGenerator) sendReminderInfinityLoop() {
	log.Println("Start sendReminderInfinityLoop")

	for {
		//telegram.sendRemind()
		telegram.sendRemindToUser(helpers.Pair{ID: telegram.id, UserName: telegram.username, RemindThisWeek: true})
		break

	}

	log.Println("End sendReminderInfinityLoop")

}

func (telegram *TelegramForGenerator) sendRemind() {
	log.Println("Start sendRemind")

	currentDayOfWeek := (int(time.Now().Weekday()) + 6) % 7

	if len(telegram.reminder[currentDayOfWeek]) != 0 && time.Now().Hour() > 17 {

		for i, val := range telegram.reminder[currentDayOfWeek] {

			if val.RemindThisWeek {

				telegram.sendRemindToUser(val)
				telegram.reminder[currentDayOfWeek][i].RemindThisWeek = false

			} else {
				telegram.reminder[currentDayOfWeek][i].RemindThisWeek = true
			}

		}

	}

	log.Println("End sendRemind")

}

func (telegram *TelegramForGenerator) sendRemindToUser(user helpers.Pair) {
	log.Println("Start sendRemindToUser")

	text := fmt.Sprintf("Привет, %s! Пришло время сдать дневник самоподготовки.", user.UserName)
	message := tgbotapi.NewMessage(user.ID, text)

	telegram.bot.Send(message)

	telegram.sendHTMLFileToUser(user.ID, user.UserName)

	log.Println("End sendRemindToUser")

}

func (telegram *TelegramForGenerator) checkUserRegistration() bool {
	log.Println("Start checkUserRegistration")

	_, err := os.Stat(pathToUserDirectories + telegram.username + ".user")

	log.Println("End checkUserRegistration")

	return err != nil

}

//GetUserValues reades user file and returns special struct
func (telegram *TelegramForGenerator) GetUserValues() *helpers.UserValues {

	log.Println("Start GetUserValues")

	//file := telegram.openFile()

	if telegram.checkUserRegistration() {

		telegram.sendMessageUserDoesNotExit()
		log.Println("End GetUserValues")
		return nil

	}

	data, _ := ioutil.ReadFile(pathToUserDirectories + telegram.username + ".user")

	start := 0
	userValues := new(helpers.UserValues)

	userValues.UserName = telegram.username

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
