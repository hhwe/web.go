package handles

import (
	"net/http"
	"webgo/models"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

var (
	user = models.User{}
)

func Login(c *gin.Context) {
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Name != "han" || user.Password != "123" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func AllUsers(c *gin.Context) {
	users, err := user.FindAll()
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, users)
}

func FindUser(c *gin.Context) {
	user, err := user.FindById(c.Param("id"))
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	c.ShouldBindJSON(&user)
	user.ID = bson.NewObjectId()
	err := user.Insert(user)
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	err := user.Update(user)
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	err := user.Delete(user)
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, user)
}
