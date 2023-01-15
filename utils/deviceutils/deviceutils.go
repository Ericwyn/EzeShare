package deviceutils

import (
	"os"
	"runtime"
)

// GetDeviceName 获取计算机名称
func GetDeviceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown-device"
	}

	return hostname
}

func GetDeviceType() string {
	return runtime.GOOS
}
