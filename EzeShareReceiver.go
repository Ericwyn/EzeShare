package main

import (
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"time"
)

func main() {
	// 开始向其他机器广播自己消息
	scanType := udpscan.UdpScanType

	scanType.StartBroadCast(10, 2*time.Second)
}
