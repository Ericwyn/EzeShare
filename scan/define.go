package scan

import "time"

type ScanTypeName string

type BroadcastMsg struct {
	Name    string
	Address string
}

type ScanCallback func(bool, []BroadcastMsg)

// ScanType 扫描实现
type ScanType struct {
	Name           ScanTypeName // 名称
	StartScan      func(ScanCallback)
	StartBroadCast func(times int, sleepDuration time.Duration)
	StopScan       func()
	StopBroadCast  func()
}
