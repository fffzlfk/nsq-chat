package models

import (
	"encoding/json"
	"log"
	"nsq-chat/config"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	channel string
	conn    *websocket.Conn
	send    chan *Message
	room    *Room
	user    *User
}

func (c *Client) read() {
	defer func() {
		c.room.leave <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(config.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(config.PongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(config.PongWait))
		return nil
	})

	for {
		msgType, msgData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error:", err)
			}
			break
		}

		if msgType != websocket.PongMessage {
			var msg Message
			if err := json.Unmarshal(msgData, &msg); err != nil {
				log.Println("Error:", err)
				break
			}

			msg.Name = c.user.Name
			msg.Channel = c.channel
			msg.User = string(c.user.ID)
			msg.TimeStamp = time.Now()

			msgData, _ = json.Marshal(msg)
		}
		if err := SendMessageToTopic(config.TopicName, msgData); err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(config.PingPeriod)
	defer func() {
		c.conn.Close()
		ticker.Stop()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if c.channel != msg.Channel {
				continue
			}
			if err := c.conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
