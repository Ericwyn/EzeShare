package apidef

const HttpApiServerPort = 23019

const ApiPathPermReq = "/api/premReq"
const ApiPathFileTransfer = "/api/fileTransfer"

// 权限请求接口

const RespCodeSuccess = 2000
const RespCodeParamError = 4001
const RespCodeServerError = 5001

type PermType string

const PermTypeOnce PermType = "Once"     // 一次发送权限
const PermTypeAlways PermType = "Always" // 永久发送权限

type PermReqRespType string

const PermReqRespAllowOnce PermReqRespType = "AllowOnce"
const PermReqRespDisAllow PermReqRespType = "DisAllow"
const PermReqRespAllowAlways PermReqRespType = "AllowAlways"

// PubResp 各端统一
type PubResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type ApiPermReq struct {
	PermType     PermType `json:"permType"` // 一次传输，或者是多次传输
	FileName     string   `json:"fileName"` // 文件名称
	FileSizeBits int64    `json:"fileSizeBits"`
	SenderName   string   `json:"senderName"`
	SenderAddr   string   `json:"senderAddr"`
	SenderPubKey string   `json:"senderPubKey"` // 发送者公钥
}

// CheckReq 参数校验, 校验不通过的时候就返回描述  string
func (apiPermReq *ApiPermReq) CheckReq() string {
	if apiPermReq.FileName == "" {
		return "file name is empty"
	}
	if apiPermReq.FileSizeBits == 0 {
		return "file size is empty"
	}
	if apiPermReq.SenderName == "" {
		return "sender name is empty"
	}
	if apiPermReq.SenderPubKey == "" {
		return "sender pub key is empty"
	}
	return ""
}

type ApiPermResp struct {
	SecToken         string          `json:"secToken"`         // receiver 使用发送者公钥加密后的 token, sender 使用这个 token 来计算 sign
	PermType         PermReqRespType `json:"permType"`         // 权限类型
	TransferId       string          `json:"transferId"`       // 文件传输的 id, 一个随机 id，sender 发送文件的时候也需要回传
	ReceiverDeviceId string          `json:"receiverDeviceId"` // 接收者的 deviceId, 方便 sender 保存信息
	//ReceiverPubKey string   // 接收者公钥
}

//// 文件传输接口
//
//type ApiFileTransferReq struct {
//	Sign       string  // md5( sender 私钥解密后的 SecToken + fileName + timeStamp)
//	TransferId string  // 来自 receiver 生成的 id
//	TimeStamp  int64   // 发送的 unix 时间戳
//	File       os.File // 文件信息
//}
//
//type ApiFileTransferResp struct {
//	Code int
//	Msg  string
//}
