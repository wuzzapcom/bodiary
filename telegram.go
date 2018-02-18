package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Syfaro/telegram-bot-api"
)

//Telegram ..
type Telegram struct {
	bot                *tgbotapi.BotAPI
	updates            <-chan tgbotapi.Update
	users              map[int64]UserData
	pathToTemplate     string
	pathToHTMLFolder   string
	pathToUsersFolders string
}

//Connect ..
func (telegram *Telegram) Connect(pathToUsersFolders string, authKey string, pathToTemplate string, pathToHTMLFolder string) error {

	log.Println("Start connect, path is " + pathToUsersFolders)

	_, err := telegram.isFolderExists(pathToUsersFolders)
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(authKey)
	if err != nil {
		return err
	}

	config := telegram.configureAPI(bot)

	updates, err := bot.GetUpdatesChan(config)
	if err != nil {
		return err
	}

	database := Database{}

	telegram.users = database.loadAllFiles(pathToUsersFolders) //make(map[int64]UserData)
	log.Println("telegram.users : ")
	log.Println(telegram.users)

	telegram.bot = bot
	telegram.updates = updates
	telegram.pathToTemplate = pathToTemplate
	telegram.pathToHTMLFolder = pathToHTMLFolder
	telegram.pathToUsersFolders = pathToUsersFolders

	return nil

}

//Start ..
func (telegram *Telegram) Start() {

	go telegram.startSendNotificationLoop()

	log.Println("Starting server")

	for update := range telegram.updates { //update.Message.Chat.ID

		if update.Message == nil {
			continue
		}

		log.Println("Message from user " + update.Message.Chat.FirstName)
		telegram.handleUpdate(update)
	}

}

func (telegram *Telegram) startSendNotificationLoop() {

	engToRuWeekdays := map[string]int{"Monday": 0, "Tuesday": 1, "Wednesday": 2, "Thursday": 3, "Friday": 4, "Saturday": 5, "Sunday": 6}

	for true {

		currentTime := time.Now()

		for key, val := range telegram.users {

			if engToRuWeekdays[currentTime.Weekday().String()] == val.DayToRemind {

				if currentTime.Hour() == val.HourToRemind {

					telegram.sendMessage(remindMessage, key)

				}

			}

		}

		time.Sleep(time.Hour)

	}

}

func (telegram *Telegram) handleUpdate(update tgbotapi.Update) {

	log.Println("update.Message.Command = " + update.Message.Command())

	if update.Message.Command() == helpCommand {

		telegram.sendMessage(helpMessage, update.Message.Chat.ID)

	} else if update.Message.Command() == createNewUserCommand {

		telegram.sendMessage(createUserMessage, update.Message.Chat.ID)

		telegram.users[update.Message.Chat.ID] = UserData{
			ID:           -1,
			DayToRemind:  0,
			HourToRemind: 0,
			Name:         "",
			Group:        "",
			StartPulse:   -1,
			EndPulse:     -1,
		}

	} else if update.Message.Command() == getDiaryCommand {

		userData, err := telegram.users[update.Message.Chat.ID]
		if !err {
			telegram.sendMessage(userNotRegisteredMessage, update.Message.Chat.ID)
			return
		}

		generator := Generator{}
		resultFile := generator.generate(userData, telegram.pathToTemplate, telegram.pathToHTMLFolder)

		telegram.sendFile(resultFile, update.Message.Chat.ID)

	} else {

		userData, err := telegram.users[update.Message.Chat.ID]
		if !err {

			telegram.sendMessage(update.Message.Text, update.Message.Chat.ID)

		} else {

			userData, errorMessage := telegram.fillUserDataWithMessage(userData, update.Message.Text)
			if errorMessage != "" {
				telegram.sendMessage(errorMessage, update.Message.Chat.ID)
			} else {
				userData.ID = update.Message.Chat.ID
				telegram.users[update.Message.Chat.ID] = userData
				telegram.sendMessage(successUserCreation, update.Message.Chat.ID)
				database := Database{}
				database.save(userData, telegram.pathToUsersFolders)
			}

		}

	}

}

func (telegram *Telegram) fillUserDataWithMessage(userData UserData, text string) (UserData, string) {

	check := telegram.checkDataCorrectness(text)
	if check != "" {
		return userData, check
	}

	data := strings.Split(text, "\n")

	userData.Name = data[0]
	userData.Group = data[1]

	userData.StartPulse, _ = strconv.Atoi(data[2])
	userData.EndPulse, _ = strconv.Atoi(data[3])

	if len(data) > 4 {

		daysOfWeek := []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}

		var i int

		for i = 0; i < 8; i++ {

			if i == 7 {
				return userData, "Такого дня недели не существует"
			} else if strings.Compare(data[4], daysOfWeek[i]) == 0 {
				userData.DayToRemind = i
				break
			}

		}

		userData.HourToRemind, _ = strconv.Atoi(data[5])
	} else {
		userData.DayToRemind = -1
		userData.HourToRemind = -1
	}

	return userData, ""

}

func (telegram *Telegram) checkDataCorrectness(text string) (errorMessage string) {

	if len(text) == 0 {
		return "Сообщение пустое"
	}

	data := strings.Split(text, "\n")

	if len(data) != minimalNumberOfFieldsInCreateUserMessage {
		return "Недостаточно данных"
	}

	n1, err := strconv.Atoi(data[2])
	if err != nil {
		return "Начальный пульс некорректен"
	}
	if n1 <= 0 {
		return "Начальный пульс некорректен"
	}

	n1, err = strconv.Atoi(data[3])
	if err != nil {
		return "Конечный пульс некорректен"
	}
	if n1 <= 0 {
		return "Конечный пульс некорректен"
	}

	if len(data) > minimalNumberOfFieldsInCreateUserMessage {

		n1, err = strconv.Atoi(data[5])
		if err != nil {
			return "Время напоминания некорректно"
		}
		if n1 < -2 || n1 > 24 {
			return "Время напоминания некорректно"
		}

	}

	return ""

}

func (telegram *Telegram) sendMessage(message string, id int64) {

	log.Println("Send message to user with message : " + message)

	telegram.bot.Send(tgbotapi.NewMessage(id, message))

}

func (telegram *Telegram) sendFile(filePath string, id int64) {

	log.Println("Send file : " + filePath)

	_, err := telegram.isFolderExists(filePath)
	if err != nil {
		log.Println("Sending file to user failed :")
		log.Println(err)
		return
	}

	telegram.bot.Send(tgbotapi.NewDocumentUpload(id, filePath))

}

func (telegram *Telegram) configureAPI(bot *tgbotapi.BotAPI) tgbotapi.UpdateConfig {

	bot.Debug = false
	config := tgbotapi.NewUpdate(0) //todo check in documentation for value
	config.Timeout = 60             //todo check in documentation for value

	return config //TODO configure this

}

func (telegram *Telegram) isFolderExists(pathToUsersFolders string) (bool, error) {

	_, err := os.Stat(pathToUsersFolders)

	if err == nil {

		return true, nil

	} else if os.IsNotExist(err) {

		log.Println("Folder doesn`t exists")

	} else {

		log.Println(err)

	}

	return false, err

}
