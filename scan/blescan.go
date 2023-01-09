package scan

const ScanTypeNameBle = "BLE"

var BleScanType = ScanType{
	Name: ScanTypeNameBle,
	ScanStart: func(callback ScanCallback) {

	},
}
