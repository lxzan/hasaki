package main

import "github.com/gin-gonic/gin"

var response = gin.H{
	"success": true,
}

func main() {
	app := gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	app.GET("/p1", func(ctx *gin.Context) {
		ctx.JSON(200, response)
	})

	app.GET("/p2", func(ctx *gin.Context) {
		ctx.JSON(400, response)
	})

	app.GET("/p4", func(ctx *gin.Context) {
		panic("test")
	})

	app.Run(":8080")
}
