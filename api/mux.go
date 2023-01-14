package api

import (
	"crypto/rand"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/gin-gonic/gin"
	"math/big"
)

// 设置 API 路由
func initAPI(router *gin.Engine) {
	router.POST(apidef.ApiPathPermReq, apiPermReq)
	router.POST(apidef.ApiPathFileTransfer, apiReceiver)
}

// NewMux 返回全局路由, 包括静态资源
func NewMux() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	initAPI(router)
	return router
}

var keyParisLen = 64

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*("

func GeneralRandomStr(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		index64 := index.Int64()
		str += letterBytes[int(index64) : int(index64)+1]
	}
	return str
}
