package telegramForGenerator

import (
	"fmt"
	"strconv"
	"wuzzapcom/bodiary/helpers"

	"github.com/Syfaro/telegram-bot-api"
)

//TODO Integrate DB and current working code

//Queue ...
type Queue struct {
	channel  chan tgbotapi.Update
	telegram *TelegramForGenerator
	users    map[int64]*helpers.UserValues
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
		case 5:
			queue.handleStateFive(update)
		}

	}
}

func (queue *Queue) handleStateZero(update tgbotapi.Update) {
	if update.Message.Command() == "CreateNewUser" {

		fmt.Printf("Test : %d, %d\n", update.Message.Chat.ID, helpers.Automaton[0][1])
		queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[0][1])
		queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите имя студента")
		queue.users[update.Message.Chat.ID] = new(helpers.UserValues)
		queue.users[update.Message.Chat.ID].ID = update.Message.Chat.ID

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

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[1][2])

	queue.users[update.Message.Chat.ID].Name = update.Message.Text

	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите группу")

}

func (queue *Queue) handleStateTwo(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[2][3])
	queue.users[update.Message.Chat.ID].Group = update.Message.Text

	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите начальный пульс")

}

func (queue *Queue) handleStateThree(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[3][4])
	pulse, _ := strconv.Atoi(update.Message.Text)
	queue.users[update.Message.Chat.ID].StartPulse = pulse

	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите конечный пульс")

}

func (queue *Queue) handleStateFour(update tgbotapi.Update) {

	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[4][5])
	pulse, _ := strconv.Atoi(update.Message.Text)
	queue.users[update.Message.Chat.ID].EndPulse = pulse
	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Введите, раз в сколько дней присылать вам напоминания")

}

func (queue *Queue) handleStateFive(update tgbotapi.Update) {
	queue.telegram.updateUserState(update.Message.Chat.ID, helpers.Automaton[5][6])
	period, _ := strconv.Atoi(update.Message.Text)
	queue.users[update.Message.Chat.ID].Period = period
	queue.users[update.Message.Chat.ID].DaysUntilSend = period
	queue.telegram.mongo.addToDB(queue.users[update.Message.Chat.ID])
	delete(queue.users, update.Message.Chat.ID)
	queue.telegram.sendQueryToUser(update.Message.Chat.ID, "Готово!")
}
