package main

import "io/ioutil"
import "encoding/json"
import "log"

/*
	TODO
	Write tests
	
	db := new(Database)

	userData := &UserData{
		ID : 0,
		DayToRemind : 2,
		HourToRemind : 3,
		Name : "wuzzapcom",
		Group : "ui9",
		StartPulse : 50,
		EndPulse : 150,
	}

	db.save(userData, "/home/wuzzapcom")
	db.load("wuzzapcom", "/home/wuzzapcom")

*/

type Database struct{}

func (database *Database) save(userData *UserData, path string){

	convertedJSON, err := json.Marshal(userData)

	log.Println(string(convertedJSON))

	if err != nil {
		log.Println("Error with converting to JSON")
		log.Println(err)
		return
	}

	ioutil.WriteFile(path + "/" + userData.Name + ".json", convertedJSON, 0644)

}

func (database *Database) load(userName string, path string) (*UserData, error) {

	data, err := ioutil.ReadFile(path + "/" + userName + ".json")
	if err != nil {
		log.Println("Error with reading file")
		log.Println(err)
		return nil, err
	}

	userData := new(UserData)

	err = json.Unmarshal(data, &userData)
	if err != nil {
		log.Println("Error with decoding JSON")
		log.Println(err)
		return nil, err
	}

	log.Println(userData)

	return userData, nil

}