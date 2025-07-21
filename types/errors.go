package types

import "github.com/gin-gonic/gin"

func ResponseWithError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gin.H{"error": err.Error()})
}
