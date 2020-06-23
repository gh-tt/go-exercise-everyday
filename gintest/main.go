package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("v1")

	v1.GET("/login", loginHandle)

	v1.GET("/info", limitRate(loginHandle))

	r.Run()
}

func loginHandle(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}

func limitRate(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {

		handler(context)
	}
}
