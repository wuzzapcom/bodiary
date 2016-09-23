package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
)

//UserValues Container for user data.
type UserValues struct {
	UserName   string
	Name       string
	Group      string
	StartPulse int
	EndPulse   int
}

type RemiderDates map[int][]Pair

type Pair struct {
	ID             int64
	UserName       string
	RemindThisWeek bool
}

//ToGOB64 function from stackoverflow
func ToGOB64(m RemiderDates) string {
	log.Println("Start ToGOB64")
	gob.Register(RemiderDates{})
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		log.Println(`failed gob Encode`, err)
	}
	log.Println("End ToGOB64")
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

//FromGOB64 function from stackoverflow
func FromGOB64(str string) RemiderDates {
	log.Println("Start FromGOB64")
	m := RemiderDates{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		log.Println(`failed gob Decode`, err)
	}
	log.Println("End FromGOB64")
	return m
}

const NumberOfExersices = 3
const NumberOfWeeks = 2
const NumberOfDays = 2

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
