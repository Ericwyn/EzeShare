package terminalui

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/api"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"github.com/Ericwyn/EzeShare/ui"
	"time"
)

func runReceiver(args ui.MainUiArgs) {
	// 开始向其他机器广播自己消息
	scanType := udpscan.UdpScanType

	// 协程执行广播, 每 2s 播发一次自己的位置
	scanType.StartBroadCastAsync(999, 2*time.Second)

	// 设置接收到文件权限请求时候的回调
	api.SetPermReqCallback(func(req apidef.ApiPermReq) apidef.PermReqRespType {
		log.I("收到来自", req.SenderName, "的文件: ", req.FileName, ", 大小: ", req.FileSizeBytes)
		log.I("是否接收? 0. 拒绝接收, 1. 接收一次, 2. 始终允许")
		var allowInput string
		fmt.Scanln(&allowInput)

		if allowInput == "0" {
			return apidef.PermReqRespDisAllow
		} else if allowInput == "1" {
			return apidef.PermReqRespAllowOnce
		} else if allowInput == "2" {
			return apidef.PermReqRespAllowAlways
		}
		return apidef.PermReqRespDisAllow
	})

	// 开启一个认证和文件接收的 api 服务器
	api.StartReceiverHttpServer()
}
