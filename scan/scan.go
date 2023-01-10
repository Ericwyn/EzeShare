package scan

//func GetScanType(typeName ScanTypeName) ScanType {
//	var scanType ScanType
//	if typeName == udpscan.ScanTypeNameUdp {
//		scanType = udpscan.UdpScanType
//	} else if typeName == blescan.ScanTypeNameBle {
//		scanType = blescan.BleScanType
//	}
//	return scanType
//}

//// StartScan 策略模式选择扫描方案
//func StartScan(typeName ScanTypeName, scanDuration time.Duration, callback ScanCallback) {
//	scanType := getScanType(typeName)
//	scanType.StartScan(callback)
//}
//
//func StopScan(typeName ScanTypeName, ) {
//
//}
//
//// StartBroadCast 开始广播消息
//func StartBroadCast(typeName ScanTypeName, broadCastTimes int, broadCastInterval int) {
//	scanType := getScanType(typeName)
//	scanType.StartBroadCast()
//}
