package scan

import "time"

type ScanTypeName string

// BroadcastMsg 广播消息，需要各端统一，此处使用小写开头驼峰命名，方便 Java JSON 格式化
type BroadcastMsg struct {
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	DeviceId   string `json:"deviceId,omitempty"`
	DeviceType string `json:"deviceType,omitempty"`
}

type ScanCallback func(bool, []BroadcastMsg)

// ScanMethod 扫描实现, BLE / UDP 扫描都是这个 struct 的一个实现
type ScanMethod struct {
	Name                ScanTypeName // 名称
	StartScanAsync      func(ScanCallback)
	StartBroadCastAsync func(times int, sleepDuration time.Duration)
	StopScan            func()
	StopBroadCast       func()
}
