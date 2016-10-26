package Telegram

import (
	"fmt"
	"wuzzapcom/bodiary/helpers"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Mongo is main struct
type Mongo struct {
	session    *mgo.Session
	collection *mgo.Collection
	telegram   *Telegram
}

func (mongo *Mongo) init(telegram *Telegram) {
	session, err := mgo.Dial("127.0.0.1")
	helpers.PanicErr(err)

	mongo.session = session

	collection := mongo.session.DB("testDB").C("Users")

	mongo.collection = collection
	mongo.telegram = telegram

}

func (mongo *Mongo) addToDB(userValues *helpers.UserValues) {

	err := mongo.collection.Insert(&userValues)
	helpers.PanicErr(err)

}

func (mongo *Mongo) getUserValuesByID(id int64) *helpers.UserValues {
	var result helpers.UserValues
	mongo.collection.Find(bson.M{"id": id}).One(&result)
	return &result
}

func (mongo *Mongo) findUsersToRemind() {

	var users []helpers.UserValues

	mongo.collection.Find(bson.M{"daysuntilsend": 0}).Iter().All(&users)

	for _, user := range users {

		mongo.telegram.sendQueryToUser(user.ID, fmt.Sprintf("%s, Hello! It`s time to print your diary", user.Name))

	}

}

func (mongo *Mongo) updateUserValues() {

	var users []helpers.UserValues
	mongo.collection.Find(bson.M{"daysuntilsend": 0}).Iter().All(&users)
	mongo.collection.RemoveAll(bson.M{"daysuntilsend": 0})

	for _, user := range users {
		user.DaysUntilSend = user.Period
		err := mongo.collection.Insert(&user)
		helpers.PanicErr(err)

	}

	change := bson.M{"$inc": bson.M{"daysuntilsend": -1}}

	query := bson.M{"daysuntilsend": bson.M{"$ne": 0}}

	mongo.collection.UpdateAll(query, change)

}
