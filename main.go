package main

import (
	"nsq-chat/controllers"
	"nsq-chat/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	room := models.NewRoom()

	router.GET("/", controllers.Login)
	router.POST("/", controllers.Login)
	router.GET("/logout", controllers.Logout)

	router.GET("/channel", controllers.ChannelList)
	router.GET("/channel/new", controllers.NewChannel)
	router.POST("/channel/new", controllers.NewChannel)
	router.GET("/channel/:id/chat", models.RoomChat(room))
	router.GET("/channel/:id/view", controllers.ChannelView)
	router.GET("/channel/:id/history", controllers.ChannelHistory)

	router.Run(":3000")
}
