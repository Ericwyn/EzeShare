package scan

const ScanTypeNameUdp = "WIFI"

var UdpScanType = ScanType{
	Name: ScanTypeNameUdp,
	ScanStart: func(callback ScanCallback) {

	},
}
