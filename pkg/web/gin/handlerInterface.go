package gin

import "github.com/gin-gonic/gin"

type Handler interface {
	Test(ctx *gin.Context)
}
