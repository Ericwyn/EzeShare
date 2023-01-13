package main

import (
	"github.com/Ericwyn/EzeShare/api"
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"time"
)

func main() {
	// 开始向其他机器广播自己消息
	scanType := udpscan.UdpScanType

	// 协程执行广播, 每 2s 播发一次自己的位置
	go scanType.StartBroadCast(999, 2*time.Second)

	// 开启一个认证和文件接收的 api 服务器
	api.StartReceiverHttpServer()
}
