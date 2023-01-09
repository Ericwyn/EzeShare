package scan

// 策略模式选择扫描方案
func StartScan(typeName string, callback ScanCallback) {
	var scanType ScanType
	if typeName == ScanTypeNameUdp {
		scanType = UdpScanType
	} else if typeName == ScanTypeNameBle {
		scanType = BleScanType
	}
	scanType.ScanStart(callback)
}
