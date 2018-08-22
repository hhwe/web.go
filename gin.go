package main

import (
	"webgo/handles"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", handles.HomePage)
	router.GET("/books", handles.AllBooks)
	router.GET("/books/:id", handles.FindBook)
	router.POST("/books", handles.CreateBook)
	router.PUT("/books", handles.UpdateBook)
	router.DELETE("/books", handles.DeleteBook)

	router.Run()
}
