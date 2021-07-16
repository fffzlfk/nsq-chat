package controllers

import (
	"nsq-chat/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

func Index(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func Login(c *fiber.Ctx) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	err := models.QueryUserByEmail(email, &user)
	if err != nil && err.Error() == "not found" {
		user.ID = bson.NewObjectId()
		user.Name = name
		user.Email = email
		user.Password = password
		user.CreatedAt = time.Now()
		if err := models.InsertUser(user); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	} else if err != nil || password != user.Password {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	cookie := &fiber.Cookie{
		Name:    "id",
		Value:   user.ID.Hex(),
		Expires: time.Now().Add(24 * time.Hour),
	}
	c.Cookie(cookie)

	return c.Render("index", nil)
}
