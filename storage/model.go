package storage

// 数据结构定义

type ConfigKey string

const ConfigKeyToken ConfigKey = "ConfigKeyToken"

// DbEzeShareConfig 数据配置表
type DbEzeShareConfig struct {
	Key   ConfigKey `xorm:"pk"`
	Value string
}

// DbEzeShareConnect 连接信息表
type DbEzeShareConnect struct {
	Name        string `xorm:"pk"`
	Type        string // 设备类型
	Token       string // 连接 token
	IsOnceToken bool   // 是否是 token
}
