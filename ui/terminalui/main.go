package terminalui

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/EzeShare/ui"
	"github.com/Ericwyn/EzeShare/utils/netutils"
	"net"
	"os"
	"strconv"
	"strings"
)

var TerminalUi = ui.UI{
	Name:      "Terminal",
	RunMainUI: runTerminalMainUi,
}

func runTerminalMainUi(args ui.MainUiArgs) {
	printLogo(args.RunMode)

	if args.IpAddr != "" {
		netutils.SetIPv4(args.IpAddr)
	} else {
		getOrSelectIPv4Addr()
	}

	if args.RunMode == ui.MainUiRunModeSender {
		runSender(args)
	} else if args.RunMode == ui.MainUiRunModeReceiver {
		runReceiver(args)
	}
}

func printLogo(mode ui.MainUiRunMode) {
	fmt.Println(
		"  _____         ____  _                    \n" +
			" | ____|_______/ ___|| |__   __ _ _ __ ___ \n" +
			" |  _| |_  / _ \\___ \\| '_ \\ / _` | '__/ _ \\\n" +
			" | |___ / /  __/___) | | | | (_| | | |  __/\n" +
			" |_____/___\\___|____/|_| |_|\\__,_|_|  \\___|\n" +
			"                                           \n" +
			"          run in mode: " + mode + "\n")
}

func getOrSelectIPv4Addr() {
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
		fmt.Println("当前设备 IP 为 ", ipList[0].String())
		netutils.SetIPv4(ipList[0].String())
		return
	}

	// 看看旧的 ip 是否在当前 ip 里面，如果在的话就直接用以前的 ip
	addr := storage.GetSelfIpAddr()
	for _, ip := range ipList {
		if addr == ip.String() {
			log.I("使用历史 IP 设置 : ", addr)
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
	fmt.Println("选择 ip 为", ipSelect.To4().String(), ", 记住 ip: ", rememberIpFlag)
	if rememberIpFlag {
		storage.SaveSelfIpAddr(ipSelect.To4().String())
	}
	netutils.SetIPv4(ipSelect.To4().String())
}
