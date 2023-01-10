package main

import (
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan"
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"sync"
)

func main() {
	// 开始扫描其他机器
	scanType := udpscan.UdpScanType

	scanType.StartScan(func(b bool, msgs []scan.BroadcastMsg) {
		if !b {
			return
		}

		log.I("receiver broadcast msg")
		for _, msg := range msgs {
			log.I("receiver ezeShareSender, name:", msg.Name, ", addr:", msg.Address)
		}
	})

	// 阻塞
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
