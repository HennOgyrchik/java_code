package ginSrvr

import "github.com/gin-gonic/gin"

type Handler interface {
	Wallet(*gin.Context)
	Wallets(*gin.Context)
	NewWallet(*gin.Context)
}
