package main

//UserData ..
type UserData struct {
	ID           int64  `json:"id"`
	DayToRemind  byte   `json:"daytoremind"`
	HourToRemind int    `json:"hourtoremind"`
	Name         string `json:"name"`
	Group        string `json:"group"`
	StartPulse   int    `json:"startpulse"`
	EndPulse     int    `json:"endpulse"`
}

//UserTemplateData ..
type UserTemplateData struct {
	Name                string
	Group               string
	FirstWeekNumber     int
	SecondWeekNumber    int
	FirstWeekExersices  [7]string
	SecondWeekExersices [7]string
}
