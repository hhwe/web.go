package handles

import (
	"net/http"
	"webgo/models"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

var (
	book = models.Book{}
)

func HomePage(c *gin.Context) {
	c.String(200, "home page")
}

func AllBooks(c *gin.Context) {
	books, err := book.FindAll()
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, books)
}

func FindBook(c *gin.Context) {
	book, err := book.FindById(c.Param("id"))
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	book.ID = bson.NewObjectId()
	book.Title = c.PostForm("title")
	book.Authors = c.PostFormArray("authors")
	book.Price = c.PostForm("price")
	err := book.Insert(book)
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, book)
}

func UpdateBook(c *gin.Context) {
	book.Title = c.PostForm("title")
	book.Authors = c.PostFormArray("authors")
	book.Price = c.PostForm("price")
	err := book.Update(book)
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	book.Title = c.PostForm("title")
	book.Authors = c.PostFormArray("authors")
	book.Price = c.PostForm("price")
	err := book.Delete(book)
	if err != nil {
		c.Err()
	}
	c.JSON(http.StatusOK, book)
}
