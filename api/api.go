package api

import (
	"github.com/Ericwyn/EzeShare/ui"
	"github.com/gin-gonic/gin"
	"sync"
)

var waitGroupMap = make(map[string]sync.WaitGroup)

func apiPermReq(ctx *gin.Context) {
	var reqBody ApiPermReq
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(200, PubResp{
			Code: RespCodeParamError,
			Msg:  "json parse error",
		})
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	ui.TerminalUi.ShowPermReqUiAsync(reqBody, func(permType PermReqRespType) {
		wg.Done()

	})

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
