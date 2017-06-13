package main

type UserData struct{

	ID int64 `json:"id"`
	DayToRemind byte `json:"daytoremind"`
	HourToRemind byte `json:"hourtoremind"`
	Name string `json:"name"`
	Group string `json:"group"`
	StartPulse int `json:"startpulse"`
	EndPulse int `json:"endpulse"`

}