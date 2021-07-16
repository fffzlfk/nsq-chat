package main

import (
	"log"
	"nsq-chat/controllers"
	"nsq-chat/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	app.Get("/", controllers.Index)

	app.Post("/login", controllers.Login)

	room := models.NewRoom()

	app.Get("/ws/:id", models.RoomChat(room))

	app.Get("/channel", controllers.ChannelList)

	app.Get("/channel/new", controllers.NewChannel)

	app.Get("/channel/:id/view", controllers.ChannelView)

	app.Get("/channel/:id/history", controllers.ChannelHistory)

	log.Fatal(app.Listen(":3000"))
}
