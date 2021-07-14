package models

import (
	"encoding/hex"
	"log"
	"nsq-chat/config"
	"nsq-chat/db"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

type Room struct {
	forward    chan *Message
	join       chan *Client
	leave      chan *Client
	clients    map[*Client]bool
	nsqReaders map[string]*NsqReader
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			if _, ok := r.clients[client]; ok {
				close(client.send)
				delete(r.clients, client)
			}
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func NewRoom() *Room {
	r := &Room{
		forward:    make(chan *Message),
		join:       make(chan *Client),
		leave:      make(chan *Client),
		clients:    make(map[*Client]bool),
		nsqReaders: make(map[string]*NsqReader),
	}

	go r.run()
	subscribeToNsq(r)
	return r
}

func RoomChat(r *Room) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomId := c.Param("id")

		upgrader := &websocket.Upgrader{
			ReadBufferSize:  config.SocketBufferSize,
			WriteBufferSize: config.SocketBufferSize,
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("ServeHTTP:", err)
			return
		}

		cookie, err := c.Request.Cookie("user_cookie")
		if err != nil {
			log.Println("Cookie Error:", err)
			return
		}
		userId := cookie.Value
		bytes, _ := hex.DecodeString(userId)
		userId = string(bytes)

		var user User
		if err := db.Mgo.DB("").C("users").FindId(bson.ObjectId(userId)).One(&user); err != nil {
			log.Println("DB Error:", err)
			return
		}

		client := &Client{
			conn:    conn,
			send:    make(chan *Message),
			room:    r,
			user:    &user,
			channel: roomId,
		}

		r.join <- client

		defer func() {
			r.leave <- client
		}()

		go client.write()
		client.read()
	}
}
