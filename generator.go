package main

import "html/template"
import "log"
import "os"
import "io/ioutil"
import "math/rand"
import "strconv"

/*
	TODO
	Write tests

	userData := UserData{
		ID : 0,
		DayToRemind : 2,
		HourToRemind : 3,
		Name : "wuzzapcom",
		Group : "ui9",
		StartPulse : 50,
		EndPulse : 150,
	}

	generator := &Generator{}

	generator.generate(userData, "/home/wuzzapcom/Coding/Golang/src/wuzzapcom/bodiary_v2.0/template.html", "/home/wuzzapcom")


*/

type Generator struct{

}

func (generator *Generator) generate(userData UserData, pathToTemplate string, pathToHTMLFolder string) (pathToFile string){

	htmlTemplate := template.New("diary")

	htmlTemplate, err := htmlTemplate.Parse(generator.openTemplate(pathToTemplate))
	if err != nil{
		log.Println("Failed parsing template")
		log.Println(err)
	}

	templateData := generator.generateUserTemplateData(userData)

	pathToFile = pathToHTMLFolder + "/" + userData.Name + ".html"

	f, err := os.Create(pathToFile)
	if err != nil{
		log.Println("Error with creating file")
		log.Println(err)
		return ""
	}

	htmlTemplate.Execute(f, templateData)

	f.Close()

	return pathToFile

}

func (generator *Generator) generateUserTemplateData(userData UserData) UserTemplateData{

	templateData := UserTemplateData{}

	templateData.Name = userData.Name
	templateData.Group = userData.Group
	templateData.FirstWeekNumber = 0
	templateData.SecondWeekNumber = 1

	for i := 0; i < 7; i++ {

		templateData.FirstWeekExersices[i] = "День отдыха"
		templateData.SecondWeekExersices[i] = "День отдыха"

	}

	for i := 0; i < 3; i++ {

		randNum := rand.Intn(7)
		for templateData.FirstWeekExersices[randNum] != "День отдыха" {
			 randNum = rand.Intn(7) 
		}
		templateData.FirstWeekExersices[randNum] = generator.generateUserExersiceDay(userData)

		randNum = rand.Intn(7)
		for templateData.SecondWeekExersices[randNum] != "День отдыха" {
			 randNum = rand.Intn(7) 
		}
		templateData.SecondWeekExersices[randNum] = generator.generateUserExersiceDay(userData)

	}

	log.Println(templateData)

	return templateData

}

func (generator *Generator) generateUserExersiceDay(userData UserData) string {

	exersices := []string{"Отжимания", "Приседания", "Пресс", "Трицепсы", "Растяжка", "Спина"}
	result := "Упражнения : \n"

	for i := 0; i < 3; i++ {

		randNum := rand.Intn(len(exersices))

		result += strconv.Itoa(i) + ") " + exersices[randNum] + "\n"

		exersices = append(exersices[ :randNum ], exersices[ randNum+1:]...)

	}

	return result

}

func (generator *Generator) openTemplate(pathToTemplate string) string{

	data, err := ioutil.ReadFile(pathToTemplate)

	if err != nil{
		log.Println("Failed reading from file")
		log.Println(err)
		panic(err)
	}

	log.Println(string(data))

	return string(data)

}