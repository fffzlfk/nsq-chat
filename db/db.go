package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

var Mgo *mgo.Session

func init() {
	mongoUrl := "mongodb://localhost/chat"

	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		panic(err)
	}

	Mgo = session

	Mgo.SetSafe(&mgo.Safe{WMode: "majority"})

	err = Mgo.DB("").C("messages").EnsureIndexKey("channel", "timestamp")
	if err != nil {
		log.Println("Error when ensure index:", err)
	}
}
