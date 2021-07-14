package models

import "time"

type Message struct {
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	Channel   string    `json:"channel"`
	User      string    `json:"user"`
	TimeStamp time.Time `json:"timestamp"`
}
