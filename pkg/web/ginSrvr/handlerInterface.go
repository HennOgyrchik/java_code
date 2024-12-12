package ginSrvr

import "github.com/gin-gonic/gin"

type Handler interface {
	Wallet(ctx *gin.Context)
	Wallets(ctx *gin.Context)
}
