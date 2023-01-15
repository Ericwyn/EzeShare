package main

import (
	"flag"
	"fmt"
	"github.com/Ericwyn/EzeShare/api"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan/udpscan"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/EzeShare/utils/netutils"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var ip = flag.String("ip", "", "set ip address")

func main() {
	flag.Parse()
	storage.InitDb(true)

	if *ip != "" {
		netutils.SetIPv4(*ip)
	} else {
		SelectIPv4()
	}

	// 开始向其他机器广播自己消息
	scanType := udpscan.UdpScanType

	// 协程执行广播, 每 2s 播发一次自己的位置
	scanType.StartBroadCastAsync(999, 2*time.Second)

	// 开启一个认证和文件接收的 api 服务器
	api.StartReceiverHttpServer()
}

func SelectIPv4() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.E(err)
	}

	ipList := make([]net.IP, 0)

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
			if ipnet.IP.To4() != nil {
				ipList = append(ipList, ipnet.IP.To4())
			}
		}
	}

	if len(ipList) == 0 {
		log.E("can't get the ip address")
		os.Exit(-1)
	}
	if len(ipList) == 1 {
		log.I("select ip : ", ipList[0].String())
		netutils.SetIPv4(ipList[0].String())
		return
	}

	// 看看旧的 ip 是否在当前 ip 里面，如果在的话就直接用以前的 ip
	addr := storage.GetSelfIpAddr()
	for _, ip := range ipList {
		if addr == ip.String() {
			log.I("select ip : ", addr)
			netutils.SetIPv4(addr)
			return
		}
	}

	fmt.Println("当前设备 IP 列表")
	fmt.Println("\t", "ip")
	for i, ip := range ipList {
		fmt.Println("["+strconv.Itoa(i)+"]", "[10"+strconv.Itoa(i)+"]", "\t", ip.String())
	}
	fmt.Println("请选择对应的 IP, [1] 为 选择第 1 个, [101] 为选择并记住第 1 个 ")
	selectInput := ""
	fmt.Scanln(&selectInput)
	if selectInput == "" {
		log.E("select ip error")
		os.Exit(-1)
	}

	rememberIpFlag := false
	if strings.HasPrefix(selectInput, "10") && selectInput != "10" {
		rememberIpFlag = true
	}
	var ipSelect net.IP
	if rememberIpFlag {
		selectInput = selectInput[2:]
	}
	selectIndex, err := strconv.Atoi(selectInput)
	if err != nil || selectIndex < 0 || selectIndex >= len(ipList) {
		log.E("select ip error")
		os.Exit(-1)
	}
	ipSelect = ipList[selectIndex]
	log.I("选择 ip 为", ipSelect.To4().String(), ", 记住 ip: ", rememberIpFlag)
	if rememberIpFlag {
		storage.SaveSelfIpAddr(ipSelect.To4().String())
	}
	netutils.SetIPv4(ipSelect.To4().String())
}
