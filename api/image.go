package api

import (
	"github.com/gin-gonic/gin"
)

func Static(ctx *gin.Context) {
	ctx.String(200, "Hello World")
}
