package http

import "os"

// 权限请求接口

type PermType string

const PermTypeOnce PermType = "Once"     // 一次发送权限
const PermTypeAlways PermType = "Always" // 永久发送权限

type PubResp struct {
	Code int
	Msg  string
	Data any
}

type ApiPermReq struct {
	PermType     PermType // 一次传输，或者是多次传输
	FileName     string   // 文件名称
	FileSizeKb   int64
	SenderName   string
	SenderPubKey string // 发送者公钥
}

type ApiPermResp struct {
	SecToken       string   // receiver 使用发送者公钥加密后的 token
	PermType       PermType // 权限类型
	ReceiverPubKey string   // 接收者公钥
}

// 文件传输接口

type ApiFileTransferReq struct {
	Signature string  // md5( sender 私钥解密后的 SecToken + fileName + timeStamp)
	TimeStamp int64   // 发送的 unix 时间戳
	File      os.File // 文件信息
}

type ApiFileTransferResp struct {
	Code int
	Msg  string
}
