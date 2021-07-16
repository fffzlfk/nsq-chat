package models

import (
	"nsq-chat/db"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UserKey int

type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Email     string        `json:"email" bson:"email"`
	Name      string        `json:"name" bson:"name"`
	Password  string        `json:"password" bson:"password"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

func InsertUser(user User) error {
	if err := db.Mgo.DB("").C("users").Insert(user); err != nil {
		return err
	}
	return nil
}

func QueryUserById(userId string, user *User) error {
	if err := db.Mgo.DB("").C("users").FindId(bson.ObjectId(userId)).One(user); err != nil {
		return err
	}
	return nil
}

func QueryUserByEmail(email string, user *User) error {
	return db.Mgo.DB("").C("users").Find(bson.M{
		"email": email,
	}).One(user)
}
