package storage

import "time"

// 数据结构定义

type ConfigKey string

const ConfigKeySelfIp ConfigKey = "ConfigKeySelfIp"

// DbEzeShareConfig 数据配置表
type DbEzeShareConfig struct {
	Key   ConfigKey `xorm:"pk"`
	Value string
}

type DeviceType string

type PermType string

const PermTypeOnce PermType = "Once"     // 一次发送权限
const PermTypeAlways PermType = "Always" // 永久发送权限

type TokenType string

const TokenTypeFromSelf TokenType = "Self"   // 自己签发给别人的
const TokenTypeFromOther TokenType = "Other" // 别人签发给自己的

// DbEzeSharePerm 连接权限信息表，保存自己签发给别人，和别人签发给自己的 Token
type DbEzeSharePerm struct {
	Id         int64
	DeviceName string     // 设备名称
	DeviceID   string     // 设备 ID
	DeviceType DeviceType // 设备类型, 类似于 Windows/Mac/Android 之类的字符串
	Token      string     // 连接 token, 已解密的 token
	PermType   PermType   // 权限类型，Once 代表是一次性权限，Always 是永久权限
	TokenType  TokenType  // 代表 Token 的来源，是自己签发给别人的，还是别人签发给自己的
}

type TransferStatus int

const TransferStatusPreSend TransferStatus = 1 // 预备发送
const TransferStatusSending TransferStatus = 2 // 正在发送
const TransferStatusFinish TransferStatus = 3  // 发送完毕

// DbEzeShareTransferMsg 传输信息表，保存文件传输记录，别人请求传输文件的时候就会记录一条信息
type DbEzeShareTransferMsg struct {
	Id             int64
	TransferId     string // 随机字符串，用来标记一次传输
	FileName       string
	FileSizeKb     int64
	FileSavePath   string         // 文件保存位置
	OnceToken      string         // 传输这个文件的 token, 为空的话，代表的是自己签发的永久 OnceToken
	TransferStatus TransferStatus // 是否已经开始传输

	FromDeviceName    string
	FromDeviceAddress string

	RequestTime time.Time
	StartTime   time.Time
	FinishTime  time.Time
}
