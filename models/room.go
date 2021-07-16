package models

import (
	"encoding/hex"
	"log"
	"nsq-chat/config"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
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

func RoomChat(r *Room) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomId := c.Params("id")

		upgrader := &websocket.FastHTTPUpgrader{
			ReadBufferSize:  config.SocketBufferSize,
			WriteBufferSize: config.SocketBufferSize,
		}
		userId := c.Cookies("id")
		bytes, _ := hex.DecodeString(userId)
		userId = string(bytes)

		err := upgrader.Upgrade(c.Context(), func(conn *websocket.Conn) {
			var user User
			if err := QueryUserById(userId, &user); err != nil {
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
		})
		if err != nil {
			log.Println("ServeHTTP:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return nil
	}
}
