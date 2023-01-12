package http

import "github.com/gin-gonic/gin"

func apiPermReq(ctx *gin.Context) {
	var reqBody ApiPermReq
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(401, PubResp{
			Code: 401,
			Msg:  "json parse error",
		})
		return
	}

	// TODO 阻塞处理 ？

	// 如果允许
	// 1. 签发 Token
	// 2. 保存文件发送请求记录
	// 如果拒绝
	// 1. 直接返回拒绝

}

func apiReceiver(ctx *gin.Context) {
	// TODO 接口请求处理
}
