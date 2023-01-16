package scan

//func GetScanType(typeName ScanTypeName) ScanMethod {
//	var scanType ScanMethod
//	if typeName == udpscan.ScanTypeNameUdp {
//		scanType = udpscan.UdpScanType
//	} else if typeName == blescan.ScanTypeNameBle {
//		scanType = blescan.BleScanType
//	}
//	return scanType
//}

//// StartScanAsync 策略模式选择扫描方案
//func StartScanAsync(typeName ScanTypeName, scanDuration time.Duration, callback ScanCallback) {
//	scanType := getScanType(typeName)
//	scanType.StartScanAsync(callback)
//}
//
//func StopScan(typeName ScanTypeName, ) {
//
//}
//
//// StartBroadCastAsync 开始广播消息
//func StartBroadCastAsync(typeName ScanTypeName, broadCastTimes int, broadCastInterval int) {
//	scanType := getScanType(typeName)
//	scanType.StartBroadCastAsync()
//}
