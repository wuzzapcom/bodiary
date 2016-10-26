package helpers

import "fmt"

//UserValues Container for user data.
type UserValues struct {
	ID            int64
	Name          string
	Group         string
	StartPulse    int
	EndPulse      int
	DaysUntilSend int
	Period        int
}

const NumberOfExersices = 3
const NumberOfWeeks = 2
const NumberOfDays = 2
const NumberOfThreads = 2

const AutomatonStates = 5
const AutomatonSignals = 7

//Automaton - transition table with structure [current state][input signal]int - new state
var Automaton = [][]int{
	[]int{0, 1, 0, 0, 0, 0, 0},
	[]int{1, 1, 2, 1, 1, 1, 1},
	[]int{2, 2, 2, 3, 2, 2, 2},
	[]int{3, 3, 3, 3, 4, 3, 3},
	[]int{4, 4, 4, 4, 4, 5, 4},
	[]int{5, 5, 5, 5, 5, 5, 0},
}

func PanicErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

const PathToAuthToken = "token.auth"

const Help = `Команды :
	/CreateNewUser - создание нового пользователя, без этого дальнейшее использование бота невозможно
	/GetDiary - получить сгенерированный файл расписания в формате html в виде вложения`

var ReasonsToMiss = []string{
	"Устал на учебе",
	"Плохое самочувствие",
	"Не было настроения",
	"Был занят",
	"Учеба",
}

var DaysOfWeek = []string{
	"Понедельник",
	"Вторник",
	"Среда",
	"Четверг",
	"Пятница",
	"Суббота",
	"Воскресенье",
}

var Exercises = []string{
	"Отжимания",
	"Приседания",
	"Пресс",
	"Спина",
	"Трицепcы",
}
