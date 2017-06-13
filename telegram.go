package main

import(
	"log"
	"os"
	"github.com/Syfaro/telegram-bot-api"
)

type Telegram struct{

	bot     *tgbotapi.BotAPI
	updates <-chan tgbotapi.Update

}

func (telegram *Telegram) Connect(pathToUsersFolders string, authKey string) error {

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

	telegram.bot = bot
	telegram.updates = updates

	return nil

}

func (telegram *Telegram) Start() {

	log.Println("Starting server")

	for update := range telegram.updates{ //update.Message.Chat.ID

		if update.Message == nil{
			continue
		}

		log.Println("Message from user " + update.Message.Chat.FirstName)

		telegram.sendMessage(update.Message.Text, update.Message.Chat.ID)
	}

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
