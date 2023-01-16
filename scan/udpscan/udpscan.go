package udpscan

import (
	"encoding/json"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan"
	"github.com/Ericwyn/EzeShare/utils/deviceutils"
	"github.com/Ericwyn/EzeShare/utils/netutils"
	"net"
	"strconv"
	"strings"
	"time"
)

const ScanTypeNameUdp scan.ScanTypeName = "WIFI"

const UdpReceiverPort = 23010

var isUdpScanIng = false
var isUdpBroadcastIng = false

var UdpScanType = scan.ScanMethod{
	Name: ScanTypeNameUdp,
	StartScanAsync: func(callback scan.ScanCallback) {
		go startUdpScan(callback)
	},
	StartBroadCastAsync: func(times int, sleepDuration time.Duration) {
		go sendUdpBroadcast(times, sleepDuration)
	},
	StopScan: func() {
		isUdpScanIng = false
	},
}

func generalMsg() string {
	msg := scan.BroadcastMsg{
		Name:       deviceutils.GetDeviceName(),
		Address:    netutils.GetIPv4().String(),
		DeviceId:   auth.GetSelfDeviceId(),
		DeviceType: deviceutils.GetDeviceType(),
	}

	marshal, err := json.Marshal(msg)
	if err != nil {
		return err.Error()
	}

	return string(marshal)
}

func parseMsg(msg string) *scan.BroadcastMsg {
	var broadcastMsg scan.BroadcastMsg
	err := json.Unmarshal([]byte(msg), &broadcastMsg)
	if err != nil {
		return nil
	}
	return &broadcastMsg
}

func sendUdpBroadcast(scanTimes int, timeDuration time.Duration) {
	if isUdpBroadcastIng {
		log.E("udp broadcast now, can't start")
		return
	}

	isUdpBroadcastIng = true
	defer func() { isUdpBroadcastIng = false }()

	broadcastMsg := generalMsg()
	if !strings.HasPrefix(broadcastMsg, "{") {
		log.E("send broad cast error, general msg fail : " + broadcastMsg)
		isUdpBroadcastIng = false
		return
	}
	// 这里设置发送者的IP地址，自己查看一下自己的IP自行设定
	laddr := net.UDPAddr{
		IP:   netutils.GetIPv4(),
		Port: UdpReceiverPort + 1000, //
	}

	// 这里设置接收者的IP地址为广播地址
	raddr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255), // net.IPv4(0, 0, 0, 0), //
		Port: UdpReceiverPort,
	}

	conn, err := net.DialUDP("udp", &laddr, &raddr)
	if err != nil {
		log.E("start broadcast fail, dail udp error")
		log.E(err.Error())
		return
	}

	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.E("close conn error while stop send broadcast")
			log.E(err.Error())
		}
	}(conn)

	for i := 0; i < scanTimes; i++ {
		if !isUdpBroadcastIng {
			log.I("stop udp broadcast")
			return
		}

		//log.D("send one msg broadcast")
		_, err := conn.Write([]byte(broadcastMsg))
		if err != nil {
			log.E(err.Error())
		}
		time.Sleep(timeDuration)
	}
}

// startUdpScan 开始 UDP 扫描，通过监听 23010 来的广播消息
func startUdpScan(callback scan.ScanCallback) {
	if isUdpScanIng {
		log.E("udp scan now, can't start")
		return
	}

	isUdpScanIng = true
	// 最后要设置 flag 为 false
	defer func() { isUdpScanIng = false }()

	addr, e := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(UdpReceiverPort))
	if e != nil {
		log.E("resolve udp add error")
		log.E(e.Error())
		return
	}

	//log.I("start receiver msg")

	// 在UDP地址上建立UDP监听,得到连接
	conn, e := net.ListenUDP("udp", addr)
	if e != nil {
		panic(e)
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.E("close conn error while stop scan")
			log.E(err.Error())
		}
	}(conn)

	// 建立缓冲区
	buffer := make([]byte, 1024)

	for {
		if !isUdpScanIng {
			log.D("stop broadcast scan")
			break
		}

		//从连接中读取内容,丢入缓冲区
		i, udpAddr, e := conn.ReadFromUDP(buffer)
		// 第一个是字节长度,第二个是udp的地址
		if e != nil {
			log.E("read from udp error")
			log.E(e.Error())
			continue
		}

		msg := string(buffer[:i])
		//log.D("scan msg: ", msg, " from ", udpAddr)

		broadcastMsg := parseMsg(msg)
		broadcastMsg.Address = udpAddr.IP.String()
		callback(true, []scan.BroadcastMsg{*broadcastMsg})
	}

}
