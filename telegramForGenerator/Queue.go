package telegramForGenerator

import (
	"fmt"
	"wuzzapcom/bodiary/helpers"

	"github.com/Syfaro/telegram-bot-api"
)

//Queue ...
type Queue struct {
	channel  chan tgbotapi.Update
	telegram *TelegramForGenerator
}

func (queue *Queue) workWithClient() {

	for update := range queue.channel {

		if update.Message == nil {
			continue
		}

		switch queue.telegram.getUserState(update.Message.Chat.ID) {

		case 0:
			queue.handleStateZero(update)
		case 1:
			queue.handleStateOne(update)
		case 2:
			queue.handleStateTwo(update)
		case 3:
			queue.handleStateThree(update)
		case 4:
			queue.handleStateFour(update)

		}

	}
}

func (queue *Queue) handleStateZero(update tgbotapi.Update) {
	if update.Message.Command() == "CreateNewUser" {

		queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[0][1])
		queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите имя студента")

		//queue.telegram.getBasicInformationFromClient(update, queue.channel)
		//userValues := queue.telegram.GetUserValues(update.Message.Chat.UserName, update.Message.Chat.ID)
		//diaryGenerator.GenerateDiary(userValues)

	} else if update.Message.Command() == "GetDiary" {

		queue.telegram.sendHTMLFileToUser(update.Message.Chat.UserName, update.Message.Chat.ID)

	} else if update.Message.Command() == "help" {

		queue.telegram.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, helpers.Help))

	} else {

		fmt.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		queue.telegram.bot.Send(msg)

	}
}

func (queue *Queue) handleStateOne(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[1][4])

	file := queue.telegram.openFile(update.Message.Chat.UserName, ".user")
	file.WriteString(update.Message.Text + "\n")
	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите группу")

}

func (queue *Queue) handleStateTwo(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[2][5])

	file := queue.telegram.openFile(update.Message.Chat.UserName, ".user")
	file.WriteString(update.Message.Text + "\n")
	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите начальный пульс")

}

func (queue *Queue) handleStateThree(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[3][6])

	file := queue.telegram.openFile(update.Message.Chat.UserName, ".user")
	file.WriteString(update.Message.Text + "\n")
	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите конечный пульс")

}

func (queue *Queue) handleStateFour(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[4][7])

	file := queue.telegram.openFile(update.Message.Chat.UserName, ".user")
	file.WriteString(update.Message.Text + "\n")
	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Готово!")

}
