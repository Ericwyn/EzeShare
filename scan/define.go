package scan

type ScanResult struct {
	Name    string
	Address string
}

type ScanCallback func(bool, []ScanResult)

// ScanType 扫描实现
type ScanType struct {
	Name      string // 名称
	ScanStart func(ScanCallback)
}
