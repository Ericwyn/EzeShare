package ui

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/scan"
)

type PermReqUiCallback func(permType apidef.PermReqRespType)
type ScanWaitUiCallback func()
type ReceiverCheckCallback func()

type ReceiverCheckUiStopScanCb func()
type ReceiverCheckUiUpdateReceiverCb func() []scan.BroadcastMsg
type ReceiverCheckUiStartSendCb func(addr string)

type UI struct {
	Name                     string
	ShowPermReqUiAsync       func(permReq apidef.ApiPermReq, callback PermReqUiCallback) // receiver 确认是否接收文件 ui
	ShowScanWaitUiAsync      func(callback ScanWaitUiCallback)                           // sender 扫描 receiver 时候展示 ui
	ShowReceiverCheckUiAsync func(c1 ReceiverCheckUiStopScanCb, c2 ReceiverCheckUiUpdateReceiverCb, c3 ReceiverCheckUiStartSendCb)
}
