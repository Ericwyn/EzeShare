package storage

import "github.com/Ericwyn/EzeShare/log"

// config 表相关操作

// SaveSelfIpAddr 保存本机 ip
func SaveSelfIpAddr(ipAddr string) {
	// 已经有配置了
	config := DbEzeShareConfig{
		Key:   ConfigKeySelfIp,
		Value: ipAddr,
	}
	if GetSelfIpAddr() != "" {
		sqlEngine.Where("key = ? ", ConfigKeySelfIp).Cols("value").Update(&config)
	} else {
		sqlEngine.InsertOne(&config)
	}
	log.I("save self ip to db:", ipAddr)
}

// GetSelfIpAddr 获取本机 ip 配置
func GetSelfIpAddr() string {
	var config DbEzeShareConfig
	sqlEngine.Where("key = ? ", ConfigKeySelfIp).Get(&config)
	return config.Value
}
