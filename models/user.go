package models

import (
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
