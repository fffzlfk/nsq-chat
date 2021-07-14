package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"nsq-chat/db"
	"nsq-chat/models"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func Login(c *gin.Context) {

	if c.Request.Method == http.MethodPost {
		name := c.PostForm("name")
		email := c.PostForm("email")
		password := c.PostForm("password")

		var user models.User
		err := db.Mgo.DB("").C("users").Find(bson.M{
			"email": email,
		}).One(&user)
		if err != nil && err.Error() == "not found" {
			user.ID = bson.NewObjectId()
			user.Name = name
			user.Email = email
			user.Password = password
			user.CreatedAt = time.Now()
			if err := db.Mgo.DB("").C("users").Insert(user); err != nil {
				fmt.Fprintln(c.Writer, err)
				return
			}

			help(c, user.ID)
			c.HTML(http.StatusOK, "index.html", nil)
			return
		}

		if password == user.Password {
			help(c, user.ID)
			c.HTML(http.StatusOK, "index.html", nil)
			return
		} else {
			fmt.Fprintln(c.Writer, errors.New("password doesn't match"))
			return
		}
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func help(c *gin.Context, id bson.ObjectId) {
	c.SetCookie("user_cookie", id.Hex(), 1000, "/", "localhost", false, true)
}
