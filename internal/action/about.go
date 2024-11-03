package action

import "github.com/gin-gonic/gin"

func About(ctx *gin.Context) {
	ctx.Writer.Write([]byte("Hello, World!"))
}
