package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		context.PureJSON(200, gin.H{
			"message": "<br>hah<br/>",
		})
	})
	r.Run()
}
