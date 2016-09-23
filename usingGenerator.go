package main

import "wuzzapcom/bodiary/telegramForGenerator"
import "log"
import "os"

func main() {

	telegram := telegramForGenerator.ConnectToTelegram()

	telegram.WorkWithClient()

	log.SetOutput(os.Stdout)

	//userValues := telegram.GetUserValues()

	//diaryGenerator.GenerateDiary(userValues)

}
