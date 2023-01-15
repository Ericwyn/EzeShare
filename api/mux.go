package api

import (
	"context"
	"crypto/rand"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"time"
)

// timeout middleware wraps the request context with a timeout
func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// 设置 API 路由
func initAPI(router *gin.Engine) {
	router.POST(apidef.ApiPathPermReq, apiPermReq)
	router.POST(apidef.ApiPathFileTransfer, apiReceiver)
}

// NewMux 返回全局路由, 包括静态资源
func NewMux() *gin.Engine {
	router := gin.Default()

	// 超长超时时间
	router.Use(gin.Logger(), timeoutMiddleware(time.Minute*30))
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
