package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"nsq-chat/db"
	"nsq-chat/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func channelList(c *gin.Context) {
	var channels []models.Channel
	db.Mgo.DB("").C("channels").Find(nil).All(&channels)
	c.HTML(http.StatusOK, "channel-list.html", channels)
}

func NewChannel(c *gin.Context) {
	data := map[string]interface{}{}

	cookie, err := c.Request.Cookie("user_cookie")
	if err != nil {
		log.Println(err)
		c.Abort()
		return
	}

	userId := cookie.Value

	if c.Request.Method == http.MethodPost {
		name := c.PostForm("name")
		channel := models.Channel{
			ID:        bson.NewObjectId(),
			Name:      name,
			CreatedBy: string(userId),
		}

		if channel.Name == "" {
			channel.Name = "No name"
		}

		if err := db.Mgo.DB("").C("channels").Insert(channel); err != nil {
			log.Println(err)
		}

		data["channel"] = channel
	}
	c.HTML(http.StatusOK, "channel-new.html", data)
}

func channelView(c *gin.Context) {
	data := map[string]interface{}{
		"Host": c.Request.Host,
	}

	var channel models.Channel
	chId := c.Param("id")
	id := bson.ObjectIdHex(chId)

	db.Mgo.DB("").C("channels").FindId(id).One(&channel)
	data["channel"] = channel

	c.HTML(http.StatusOK, "channel-view.html", data)
}

func channelHistory(c *gin.Context) {
	const limit = 10
	result := make([]models.Message, limit)

	err := db.Mgo.DB("").C("messages").Find(
		bson.M{"channel": c.Param("id")},
	).Sort("-timestamp").Limit(limit).All(&result)

	if err != nil {
		log.Print(err)
	}

	if err := json.NewEncoder(c.Writer).Encode(result); err != nil {
		log.Println(err)
	}
}

func MustAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Request.Cookie("user_cookie")
		if err != nil {
			log.Println(err)
			c.Writer.Write([]byte("Error"))
			c.Abort()
			return
		}
		handler(c)
	}
}

var (
	ChannelList    = MustAuth(channelList)
	ChannelHistory = MustAuth(channelHistory)
	ChannelView    = MustAuth(channelView)
)
