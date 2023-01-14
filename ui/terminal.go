package ui

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/log"
)

var TerminalUi UI = UI{
	Name: "Terminal",
	ShowPermReqUiAsync: func(req apidef.ApiPermReq, callback PermReqUiCallback) {
		go terminalShowPermReqUi(req, callback)
	},
}

func terminalShowPermReqUi(req apidef.ApiPermReq, callback PermReqUiCallback) {
	// TODO 入队实现，一次只处理一个请求

	log.I("收到来自", req.SenderName, "的文件: ", req.FileName, ", 大小: ", req.FileSizeKb)
	log.I("是否接收? 0. 拒绝接收, 1. 接收一次, 2. 始终允许")
	var allowInput string
	fmt.Scanln(&allowInput)

	if allowInput == "0" {
		callback(apidef.PermReqRespDisAllow)
	} else if allowInput == "1" {
		callback(apidef.PermReqRespAllowOnce)
	} else if allowInput == "2" {
		callback(apidef.PermReqRespAllowAlways)
	}
}
