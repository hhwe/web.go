package main

import (
	"webgo/handles"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", handles.HomePage)
	router.POST("/login", handles.Login)

	v1 := router.Group("api/v1")
	{
		v1.GET("/users", handles.AllUsers)
		v1.GET("/users/:id", handles.FindUser)
		v1.POST("/users", handles.CreateUser)
		v1.PATCH("/users", handles.UpdateUser)
		v1.DELETE("/users", handles.DeleteUser)

		v1.GET("/books", handles.AllBooks)
		v1.GET("/books/:id", handles.FindBook)
		v1.POST("/books", handles.CreateBook)
		v1.PATCH("/books", handles.UpdateBook)
		v1.DELETE("/books", handles.DeleteBook)
	}

	router.Run()
}
