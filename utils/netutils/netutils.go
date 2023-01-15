package netutils

import (
	"github.com/Ericwyn/EzeShare/log"
	"net"
	"strings"
)

var ipAddrCache net.IP = nil

func SetIPv4(ipAddress string) {
	if !strings.HasSuffix(ipAddress, "/24") ||
		!strings.HasSuffix(ipAddress, "/64") {
		ipAddress = ipAddress + "/24"
	}
	ip, _, err := net.ParseCIDR(ipAddress)
	if err != nil {
		log.E("Set IP Error ", ipAddress)
		log.E(err)
		panic(err)
	}

	ipAddrCache = ip
}

func GetIPv4() net.IP {
	if ipAddrCache != nil {
		return ipAddrCache
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.E(err)
	}
	var ip string = "localhost"
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				log.D("获取到 IP 地址: " + ip)
				return ipnet.IP
			}
		}
	}
	return net.IPv4(0, 0, 0, 0)
}
