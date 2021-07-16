package models

import (
	"nsq-chat/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	Channel   string    `json:"channel"`
	User      string    `json:"user"`
	TimeStamp time.Time `json:"timestamp"`
}

func QueryMessageByChannelId(channelId string, limit int) ([]Message, error) {
	messages := make([]Message, limit)
	if err := db.Mgo.DB("").C("messages").Find(
		bson.M{"channel": channelId},
	).Sort("-timestamp").Limit(limit).All(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}
