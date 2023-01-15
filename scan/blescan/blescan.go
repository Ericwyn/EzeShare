package blescan

import "github.com/Ericwyn/EzeShare/scan"

const ScanTypeNameBle scan.ScanTypeName = "BLE"

var BleScanType = scan.ScanType{
	Name: ScanTypeNameBle,
	StartScanAsync: func(callback scan.ScanCallback) {

	},
}
