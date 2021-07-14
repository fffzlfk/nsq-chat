package models

import "gopkg.in/mgo.v2/bson"

type Channel struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	CreatedBy string        `json:"created_by" bson:"created_by"`
}
