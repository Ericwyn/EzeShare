package main

import (
	"flag"
	"fmt"
	"github.com/Ericwyn/EzeShare/api/apiclient"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan"
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"github.com/Ericwyn/GoTools/file"
	"os"
	"strconv"
	"sync"
)

var filePath = flag.String("f", "", "file path of the file which will send to others")

func main() {
	flag.Parse()
	if *filePath == "" {
		log.E("file for send is empty")
		return
	}

	// 开始扫描其他机器
	scanType := udpscan.UdpScanType

	receiverList := make([]scan.BroadcastMsg, 0)
	receiverAddressMap := make(map[string]int)
	scanType.StartScanAsync(func(b bool, msgs []scan.BroadcastMsg) {
		if !b {
			return
		}

		//log.I("receiver broadcast msg")
		for _, msg := range msgs {
			//log.I("receiver ezeShareSender, name:", msg.Name, ", addr:", msg.Address)
			if receiverAddressMap[msg.Address] != 1 {
				receiverList = append(receiverList, msg)
				receiverAddressMap[msg.Address] = 1
				printAllReceivers(receiverList)
			}
		}
	})

	var selectIdx = ""
	go func() {
		fmt.Scanln(&selectIdx)
		if selectIdx == "" {
			log.E("select receiver id error")
			os.Exit(-1)
		}
		i, err := strconv.Atoi(selectIdx)
		if err != nil || len(receiverList) <= i || i < 0 {
			log.E("receiver id error, " + selectIdx)
			os.Exit(-1)
		}
		scanType.StopScan()
		msg := receiverList[i]

		apiclient.DoPermRequest(msg.Address, file.OpenFile(*filePath), apidef.PermTypeAlways)
	}()

	// 阻塞
	var wg sync.WaitGroup
	//wgDone := false

	//// 展示扫描结果 + 结果选择 + 发送调用
	//// 需要将扫描结果实时传进去，并且需要支持停止扫描
	//// 参数是 3 个 callback
	//// 1. 触发停止扫描的 callback, ui 那边用户输入停止扫描的时候触发，这时候会回来停掉 scan
	//// 2. 更新 receiver callback 到 ui, ui 那边可以及时更新当前显示的内容
	//// 3. 触发发送的 callback, ui 那边用户选择完 receiver 之后, 开始发送
	//stopScanCallback := func() {
	//	scanType.StopScan() // 停止扫描
	//}
	//receiverUpdateCallback := func() []scan.BroadcastMsg {
	//	return receiverList
	//}
	//startSendCallback := func(addr string) {
	//	apiclient.DoPermRequest(addr, file.OpenFile(*filePath), apidef.PermTypeAlways)
	//}
	//
	//ui.TerminalUi.ShowReceiverCheckUiAsync(stopScanCallback, receiverUpdateCallback, startSendCallback)

	//go func() {
	//	// 展示
	//	ui.TerminalUi.ShowReceiverCheckUiAsync()
	//	// 展示一个 ui
	//	ui.TerminalUi.ShowScanWaitUiAsync(func() {
	//		errutils.Try(func() {
	//			if !wgDone {
	//				wg.Done()
	//				wgDone = true
	//			}
	//		}, func(i interface{}) {
	//			log.I("[ShowScanWaitUiAsync] wait group done panic, ", i)
	//		})
	//		scanType.StopScan() // 停止扫描
	//	})
	//}()

	wg.Add(1)
	wg.Wait()
}

func printAllReceivers(receiver []scan.BroadcastMsg) {
	for i, msg := range receiver {
		fmt.Println("\n\n\n\n\n\n\n\n\n")
		fmt.Println("当前 receiver 列表如下: ")
		fmt.Println("\taddress       \t", "name")
		fmt.Println("["+strconv.Itoa(i)+"]\t", msg.Address, "\t", msg.Name)
		fmt.Println("-----------------------")
		fmt.Println("输入编号并回车, 选择具体 receiver")
	}
}
