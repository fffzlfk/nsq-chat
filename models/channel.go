package models

import (
	"nsq-chat/db"

	"gopkg.in/mgo.v2/bson"
)

type Channel struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	CreatedBy string        `json:"created_by" bson:"created_by"`
}

func QueryAllChannels(channels *([]Channel)) error {
	return db.Mgo.DB("").C("channels").Find(nil).All(channels)
}

func InsertChannel(channel Channel) error {
	_, err := db.Mgo.DB("").C("channels").Upsert(bson.M{
		"name": channel.Name,
	}, channel)
	return err
}

func QueryChannelById(id bson.ObjectId, channel *Channel) error {
	return db.Mgo.DB("").C("channels").FindId(id).One(channel)
}
