package main

import "io/ioutil"
import "encoding/json"
import "log"
import "strings"
import "strconv"

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

//Database ..
type Database struct{}

func (database *Database) save(userData UserData, path string){

	convertedJSON, err := json.Marshal(userData)

	log.Println(string(convertedJSON))

	if err != nil {
		log.Println("Error with converting to JSON")
		log.Println(err)
		return
	}

	ioutil.WriteFile(path + "/" + strconv.FormatInt(userData.ID, 10) + ".json", convertedJSON, 0644)

}

func (database *Database) load(userName string, path string) (UserData, error) {

	userData := UserData{}

	data, err := ioutil.ReadFile(path + "/" + userName)
	if err != nil {
		log.Println("Error with reading file")
		log.Println(err)
		return userData, err
	}

	err = json.Unmarshal(data, &userData)
	if err != nil {
		log.Println("Error with decoding JSON")
		log.Println(err)
		return userData, err
	}

	log.Println(userData)

	return userData, nil

}

func (database *Database) loadAllFiles(path string) map[int64]UserData{

	result := make(map[int64]UserData)

	files, err := ioutil.ReadDir(path)
	if err != nil{
		log.Println("Failed reading dir")
		log.Println(err)
		return nil
	}

	for _, file := range files {

		if strings.Contains(file.Name(), ".json") {

			userData, err := database.load(file.Name(), path)
			if err != nil{
				continue
			}
			result[userData.ID] = userData

		}

	}

	return result


}