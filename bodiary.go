package main

import "os"
import "fmt"
import "io/ioutil"
import "log"

func main(){

	if len(os.Args) != 5{

		fmt.Println("Set folder with user metadata or auth file or path to HTML template or path to user`s htmls")
		return

	}

	telegram := new(Telegram)

	err := telegram.Connect(os.Args[1], openAuthFile(os.Args[2]), os.Args[3], os.Args[4])

	if err != nil {
		log.Println("Connect problem")
		log.Println(err)
		return
	}

	telegram.Start()

}

func openAuthFile(pathToAuthFile string) string{

	data, err := ioutil.ReadFile(pathToAuthFile)

	if err != nil{
		log.Println("Failed reading from file")
		log.Println(err)
		panic(err)
	}

	log.Println(string(data))

	return string(data)

}