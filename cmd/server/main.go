package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("Hello, World!"))
	})
	r.Run(":8080")
	fmt.Println()
}
