package ui

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/api"
	"github.com/Ericwyn/EzeShare/log"
)

var TerminalUi UI = UI{
	Name: "Terminal",
	ShowPermReqUiAsync: func(req api.ApiPermReq, callback PermReqUiCallback) {
		go terminalShowPermReqUi(req, callback)
	},
}

func terminalShowPermReqUi(req api.ApiPermReq, callback PermReqUiCallback) {
	// TODO 入队实现，一次只处理一个请求

	log.I("收到来自", req.SenderName, "的文件: ", req.FileName, ", 大小: ", req.FileSizeKb)
	log.I("是否接收? 0. 拒绝接收, 1. 接收一次, 2. 始终允许")
	var allowInput string
	fmt.Scanln(&allowInput)

	if allowInput == "0" {
		callback(api.PermReqRespDisAllow)
	} else if allowInput == "1" {
		callback(api.PermReqRespAllowOnce)
	} else if allowInput == "2" {
		callback(api.PermReqRespAllowAlways)
	}
}
