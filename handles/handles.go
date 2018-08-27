package handles

import "github.com/gin-gonic/gin"

func HomePage(c *gin.Context) {
	c.String(200, "home page")
}
