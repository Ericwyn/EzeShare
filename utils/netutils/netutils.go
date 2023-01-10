package netutils

import (
	"github.com/Ericwyn/EzeShare/log"
	"net"
	"os"
)

func GetIPv4() net.IP {
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

// GetDeviceName 获取计算机名称
func GetDeviceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-device"
	}

	return hostname
}
