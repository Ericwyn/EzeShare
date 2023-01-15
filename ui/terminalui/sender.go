package terminalui

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/api/apiclient"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan"
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"github.com/Ericwyn/EzeShare/ui"
	"github.com/Ericwyn/GoTools/file"
	"os"
	"strconv"
)

func runSender(args ui.MainUiArgs) {
	if args.SenderFile == "" {
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
		for _, msg := range msgs {
			if receiverAddressMap[msg.Address] != 1 {
				receiverList = append(receiverList, msg)
				receiverAddressMap[msg.Address] = 1
				printAllReceivers(receiverList)
			}
		}
	})

	var selectIdx = ""

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

	fileForUpload := file.OpenFile(args.SenderFile)

	apiclient.DoPermRequest(msg.Address,
		fileForUpload,
		apidef.PermTypeAlways,
		func(per int) {
			// 上传文件百分比回调
			printUploadProcess(fileForUpload.Name(), per)
		})
}

func printAllReceivers(receiver []scan.BroadcastMsg) {
	fmt.Println("\n\n\n\n\n\n\n\n\n")
	fmt.Println("当前 receiver 列表如下: ")
	fmt.Println("\t address", "\t\t", "name\t\t\t", "deviceId")
	for i, msg := range receiver {
		fmt.Println("["+strconv.Itoa(i)+"]\t", msg.Address, "\t", msg.Name, "\t", msg.DeviceId)
	}
	fmt.Println("-----------------------")
	fmt.Println("输入编号并回车, 选择具体 receiver")
}

func printUploadProcess(fileName string, per int) {
	process := "["
	for i := 0; i <= 100; i += 5 {
		if i < per {
			process += "="
		} else {
			process += " "
		}
	}
	process += "]"
	log.I("上传 ", fileName, ", 进度: ", process)
}
