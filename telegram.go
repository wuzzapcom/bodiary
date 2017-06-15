package main

import(
	"log"
	"os"
	"github.com/Syfaro/telegram-bot-api"
	"strconv"
	"strings"
)

type Telegram struct{

	bot     *tgbotapi.BotAPI
	updates <-chan tgbotapi.Update
	users map[int64]UserData
	pathToTemplate string
	pathToHTMLFolder string
	pathToUsersFolders string

}

func (telegram *Telegram) Connect(pathToUsersFolders string, authKey string, pathToTemplate string, pathToHTMLFolder string) error {

	log.Println("Start connect, path is " + pathToUsersFolders)

	_, err := telegram.isFolderExists(pathToUsersFolders)
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(authKey)
	if err != nil{
		return err
	}

	config := telegram.configureAPI(bot)

	updates, err := bot.GetUpdatesChan(config)
	if err != nil{
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

func (telegram *Telegram) Start() {

	log.Println("Starting server")

	for update := range telegram.updates{ //update.Message.Chat.ID

		if update.Message == nil{
			continue
		}

		log.Println("Message from user " + update.Message.Chat.FirstName)
		telegram.handleUpdate(update)
	}

}

func (telegram *Telegram) handleUpdate(update tgbotapi.Update){

	log.Println("update.Message.Command = " + update.Message.Command())

	if update.Message.Command() == "help"{

		telegram.sendMessage(HELP_MESSAGE, update.Message.Chat.ID)

	}else if update.Message.Command() == "createNewUser"{

		telegram.sendMessage(CREATE_USER_MESSAGE, update.Message.Chat.ID)

		telegram.users[update.Message.Chat.ID] = UserData{
			ID : -1,
			DayToRemind : 0,
			HourToRemind : 0,
			Name : "",
			Group : "",
			StartPulse : -1,
			EndPulse : -1,
		}

	}else if (update.Message.Command() == "getDiary"){

		userData, err := telegram.users[update.Message.Chat.ID]
		if !err {
			telegram.sendMessage(USER_NOT_REGISTERED_MESSAGE, update.Message.Chat.ID)
			return
		}

		generator := Generator{}
		resultFile := generator.generate(userData, telegram.pathToTemplate, telegram.pathToHTMLFolder)

		telegram.sendFile(resultFile, update.Message.Chat.ID)

	}else{

		userData, err := telegram.users[update.Message.Chat.ID]
		if !err{

			telegram.sendMessage(update.Message.Text, update.Message.Chat.ID)

		}else{

			userData, errorMessage := telegram.fillUserDataWithMessage(userData, update.Message.Text)
			if errorMessage != ""{
				telegram.sendMessage(errorMessage, update.Message.Chat.ID)
			}else{
				userData.ID = update.Message.Chat.ID
				telegram.users[update.Message.Chat.ID] = userData
				telegram.sendMessage(SUCCESS_USER_CREATION, update.Message.Chat.ID)
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

	daysOfWeek := []string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота", "Воскресенье"}

	var i byte

	for i = 0; i < 7; i++ {

		if data[4] == daysOfWeek[i]{
			userData.DayToRemind = i
		}

	}

	if i != 7 {
		return userData, "Такого дня недели не существует"
	}

	userData.Name = data[0]
	userData.Group = data[1]

	userData.StartPulse, _ = strconv.Atoi(data[2])
	userData.EndPulse, _ = strconv.Atoi(data[3])

	userData.HourToRemind, _ = strconv.Atoi(data[5])

	return userData, ""

}

func (telegram *Telegram) checkDataCorrectness(text string) (errorMessage string){

	if len(text) == 0 {
		return "Сообщение пустое"
	}

	data := strings.Split(text, "\n")
	
	if len(data) != NUMBER_OF_FIELDS_IN_CREATE_USER_MESSAGE {
		return "Недостаточно данных"
	}

	n1, err := strconv.Atoi(data[2])
	if err != nil{
		return "Начальный пульс некорректен"
	} 
	if n1 <= 0 {
		return "Начальный пульс некорректен"
	}

	n1, err = strconv.Atoi(data[3])
	if err != nil{
		return "Конечный пульс некорректен"
	} 
	if n1 <= 0 {
		return "Конечный пульс некорректен"
	}

	n1, err = strconv.Atoi(data[5])
	if err != nil{
		return "Время напоминания некорректно"
	}
	if n1 < 0 || n1 > 24{
		return "Время напоминания некорректно"
	}

	return ""

}

func (telegram *Telegram) sendMessage(message string, id int64){

	log.Println("Send message to user with message : " + message)

	telegram.bot.Send(tgbotapi.NewMessage(id, message))

}

func (telegram *Telegram) sendFile(filePath string, id int64){

	log.Println("Send file : " + filePath)

	_, err := telegram.isFolderExists(filePath)
	if (err != nil){
		log.Println("Sending file to user failed :")
		log.Println(err)
		return
	}

	telegram.bot.Send(tgbotapi.NewDocumentUpload(id, filePath))


}

func (telegram *Telegram) configureAPI(bot *tgbotapi.BotAPI) tgbotapi.UpdateConfig{

	bot.Debug = false
	config := tgbotapi.NewUpdate(0)//todo check in documentation for value
	config.Timeout = 60 //todo check in documentation for value

	return config //TODO configure this

}

func (telegram *Telegram) isFolderExists(pathToUsersFolders string) (bool, error) {

	_, err := os.Stat(pathToUsersFolders)

	if err == nil {

		return true, nil

	} else if os.IsNotExist(err) {

		log.Println("Folder doesn`t exists")

	} else{

		log.Println(err)

	}

	return false, err

}
