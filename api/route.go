package api

import (
	"github.com/gin-gonic/gin"
)

func List(ctx *gin.Context) {
	ctx.String(200, "Hello World")
}
