package controllers

import (
	"encoding/json"
	"nsq-chat/models"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

func ChannelList(c *fiber.Ctx) error {
	var channels []models.Channel
	if err := models.QueryAllChannels(&channels); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("channel-list", channels)
}

func NewChannel(c *fiber.Ctx) error {
	userId := c.Cookies("id")
	if userId == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	data := map[string]interface{}{}

	name := c.Query("name")
	if name != "" {
		channel := models.Channel{
			ID:        bson.NewObjectId(),
			Name:      name,
			CreatedBy: string(userId),
		}

		if channel.Name == "" {
			channel.Name = "No name"
		}

		if err := models.InsertChannel(channel); err != nil {
			c.SendString("Alreadly exists")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		data["channel"] = channel
	}
	return c.Render("channel-new", data)
}

func ChannelView(c *fiber.Ctx) error {
	data := fiber.Map{
		"Host": string(c.Request().Host()),
	}

	var channel models.Channel
	chId := c.Params("id")
	id := bson.ObjectIdHex(chId)

	if err := models.QueryChannelById(id, &channel); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	data["channel"] = channel

	return c.Render("channel-view", data)
}

func ChannelHistory(c *fiber.Ctx) error {
	channelId := c.Params("id")

	const limit = 10

	res, err := models.QueryMessageByChannelId(channelId, limit)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return json.NewEncoder(c.Context().Response.BodyWriter()).Encode(res)
}
