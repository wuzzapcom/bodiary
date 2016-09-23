package diaryGenerator

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"wuzzapcom/diary/helpers"
)

type htmlGenerator struct {
	htmlFile   *os.File
	pickedDays DaysMap
	userValues *helpers.UserValues
}

func (generator *htmlGenerator) generateSkeletonBegin() {
	log.Println("Start generateSkeletonBegin")
	defer log.Println("End generateSkeletonBegin")

	generator.htmlFile.WriteString("<html>\n")
	generator.htmlFile.WriteString(" <head>\n")
	generator.htmlFile.WriteString("  <meta charset=\"utf-8\">\n")
	generator.htmlFile.WriteString("  <title>Дневник самоподготовки</title>\n")
	generator.htmlFile.WriteString("</head>\n")
	generator.htmlFile.WriteString("<body>\n")
	generator.htmlFile.WriteString("<table border=\"1\">\n")
	generator.htmlFile.WriteString(fmt.Sprintf("  <caption>Дневник самоподготовки : %s, %s</caption>\n", generator.userValues.Name, generator.userValues.Group))
	generator.htmlFile.WriteString("<tr>\n")
	generator.htmlFile.WriteString("<th>День недели</th>\n")
	generator.htmlFile.WriteString("<th>Неделя 1</th>\n")
	generator.htmlFile.WriteString("<th>Неделя 2</th>\n")
	generator.htmlFile.WriteString("</tr>")

}

func (generator *htmlGenerator) generateSkeletonBody() {
	log.Println("Start generateSkeletonBody")
	defer log.Println("End generateSkeletonBody")

	for i := 0; i < 7; i++ {

		generator.htmlFile.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", helpers.DaysOfWeek[i], generator.pickedDays[i][0], generator.pickedDays[i][1]))

	}

}

func (generator *htmlGenerator) generateSkeletonEnd() {
	log.Println("Start generateSkeletonEnd")
	defer log.Println("End generateSkeletonEnd")

	generator.htmlFile.WriteString("  </table>")
	generator.htmlFile.WriteString(" </body>")
	generator.htmlFile.WriteString("</html>")

}

func (generator *htmlGenerator) generateRows() {
	log.Println("Start generateRows")
	defer log.Println("End generateRows")

	for i := 0; i < 7; i++ {

		generator.pickedDays.addDay(i)

	}

	for i := 0; i < helpers.NumberOfWeeks; i++ {

		for j := 0; j < helpers.NumberOfDays; j++ {

			day := generator.generateDayOfWeek()
			for i := 0; i < 2 && generator.checkForTwoDaysInRow(day); i++ {

				day = generator.generateDayOfWeek()

			}

			if rand.Intn(10) == 0 && day != 6 {

				generator.pickedDays.addMessage(day, i, generator.generateReasonToMiss())

				fmt.Println("Generated reason to miss : ", generator.pickedDays[day][i])

			} else {

				generator.pickedDays.addMessage(day, i, generator.generateTrainingResult())

			}

		}

	}

}

func (generator *htmlGenerator) generateTrainingResult() string {
	log.Println("Start generateTrainingResult")
	defer log.Println("End generateTrainingResult")

	pulseStart := generator.userValues.StartPulse

	pulseEnd := generator.userValues.EndPulse

	delta1 := rand.Intn(10) - 5

	delta2 := rand.Intn(10) - 5

	exercisesInTraining := generator.generateExercises()

	return fmt.Sprintf("Упражнения : <br>1) %s<br>2) %s<br>3) %s.<br>Пульс в начале : %d, Пульс в конце : %d<br>Самочувствие : хорошее",
		exercisesInTraining[0], exercisesInTraining[1], exercisesInTraining[2], pulseStart+delta1, pulseEnd+delta2)

}

func (generator *htmlGenerator) generateReasonToMiss() string {
	log.Println("Start generateReasonToMiss")
	defer log.Println("End generateReasonToMiss")

	return helpers.ReasonsToMiss[rand.Intn(len(helpers.ReasonsToMiss))]

}

func (generator *htmlGenerator) generateDayOfWeek() int {
	log.Println("Start generateDayOfWeek")
	defer log.Println("End generateDayOfWeek")

	return rand.Intn(7)

}

func (generator *htmlGenerator) generateExercises() []string {
	log.Println("Start generateExercises")
	defer log.Println("End generateExercises")

	result := make([]string, helpers.NumberOfExersices)
	resultInt := make([]int, helpers.NumberOfExersices)

	resultInt[0] = rand.Intn(len(helpers.Exercises))
	for i := 1; i < helpers.NumberOfExersices; i++ {

		resultInt[i] = rand.Intn(len(helpers.Exercises))

		for j := 0; j < i; j++ {

			if resultInt[i] == resultInt[j] {

				i--
				break

			}

		}

	}

	for i, exer := range resultInt {
		result[i] = helpers.Exercises[exer]
	}

	return result

}

func getTimeInMilliseconds() int64 {
	log.Println("Start getTimeInMilliseconds")
	defer log.Println("End getTimeInMilliseconds")
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (generator *htmlGenerator) checkForTwoDaysInRow(day int) bool {

	log.Println("Start checkForTwoDaysInRow")
	defer log.Println("End checkForTwoDaysInRow")

	if day == 0 && generator.pickedDays.contains(day+1) {

		return true

	} else if day == 6 && generator.pickedDays.contains(day-1) {

		return true

	} else if generator.pickedDays.contains(day-1) ||
		generator.pickedDays.contains(day+1) {

		return true

	}

	return false

}

//GenerateDiary creates a html table.
func GenerateDiary(userValues *helpers.UserValues) {
	log.Println("Start GenerateDiary")
	rand.Seed(getTimeInMilliseconds())
	file, _ := os.Create(fmt.Sprintf("%s.html", userValues.UserName))
	generator := htmlGenerator{file, make(DaysMap), userValues}

	generator.generateRows()

	generator.generateSkeletonBegin()
	generator.generateSkeletonBody()
	generator.generateSkeletonEnd()
	log.Println("End GenerateDiary")
}
