package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	pathToUsersFolders := flag.String(pathToUsersFoldersFlag, "", pathToAuthFileHelp)

	pathToAuthFile := flag.String(pathToAuthFileFlag, "auth.private", pathToAuthFileHelp)

	pathToHTMLFile := flag.String(pathToHTMLTemplateFlag, "template.html", pathToHTMLTempleteHelp)

	pathToGeneratedFolder := flag.String(pathToGeneratedFilesFlag, "", pathToGeneratedFilesHelp)

	flag.Parse()

	if len(os.Args) != 5 {

		fmt.Println("Set folder with user metadata or auth file or path to HTML template or path to user`s htmls")
		return

	}

	telegram := new(Telegram)

	err := telegram.Connect(*pathToUsersFolders, openAuthFile(*pathToAuthFile), *pathToHTMLFile, *pathToGeneratedFolder)

	if err != nil {
		log.Println("Connect problem")
		log.Println(err)
		return
	}

	telegram.Start()

}

func openAuthFile(pathToAuthFile string) string {

	data, err := ioutil.ReadFile(pathToAuthFile)

	if err != nil {
		log.Println("Failed reading from file")
		log.Println(err)
		panic(err)
	}

	log.Println(string(data))

	return string(data)

}
