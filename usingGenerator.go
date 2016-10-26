package main

import "wuzzapcom/bodiary/Telegram"

func main() {

	telegram := Telegram.ConnectToTelegram()

	telegram.WorkWithClient()

}
