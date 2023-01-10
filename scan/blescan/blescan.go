package blescan

import "github.com/Ericwyn/EzeShare/scan"

const ScanTypeNameBle scan.ScanTypeName = "BLE"

var BleScanType = scan.ScanType{
	Name: ScanTypeNameBle,
	StartScan: func(callback scan.ScanCallback) {

	},
}
