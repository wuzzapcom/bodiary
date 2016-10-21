package telegramForGenerator

import (
	"wuzzapcom/bodiary/helpers"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Mongo is main struct
type Mongo struct {
	session    *mgo.Session
	collection *mgo.Collection
	telegram   *TelegramForGenerator
}

func (mongo *Mongo) init(telegram *TelegramForGenerator) {
	session, err := mgo.Dial("127.0.0.1")

	if err != nil {
		panic(err)
	}

	mongo.session = session

	collection := mongo.session.DB("testDB").C("Users")

	mongo.collection = collection
	mongo.telegram = telegram

}

func (mongo *Mongo) addToDB(userValues *helpers.UserValues) {

	err := mongo.collection.Insert(&userValues)

	if err != nil {
		panic(err)
	}

}

func (mongo *Mongo) findUsersToRemind() {

	var users []helpers.UserValues

	mongo.collection.Find(bson.M{"DaysUntilSend": 0}).Iter().All(&users)

	for _, user := range users {

		//mongo.telegram.sendQuery(user.)    //TODO FINISH THIS STRING

	}

}

func (mongo *Mongo) updateUserValues() {

	query := bson.M{"DaysUntilSend": bson.M{"$ne": 0}}

	change := bson.M{"$set": "$dec"}

	mongo.collection.UpdateAll(query, change)

	var users []helpers.UserValues
	mongo.collection.Find(bson.M{"DaysUntilSend": 0}).Iter().All(&users)
	mongo.collection.RemoveAll(bson.M{"DaysUntilSend": 0})

	for _, user := range users {
		user.DaysUntilSend = user.Period
		err := mongo.collection.Insert(&user)
		if err != nil {
			panic(err)
		}
	}
}
