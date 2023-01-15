package storage

import "github.com/Ericwyn/EzeShare/log"

// perm 表相关操作

// SaveSelfToken 保存 OnceToken
func SaveSelfToken(token string, deviceId string) {
	tokenSelf := DbEzeSharePerm{
		DeviceName: "Self",
		DeviceType: "Self",
		DeviceId:   deviceId,
		Token:      token,
		PermType:   PermTypeAlways,
		TokenType:  TokenTypeFromSelf,
	}

	_, err := sqlEngine.InsertOne(tokenSelf)
	if err != nil {
		log.E("save token error")
		panic(err)
	}
}

// GetSelfPermFromDB 获取 token，存在的话为 bool, string
func GetSelfPermFromDB() (bool, DbEzeSharePerm) {
	var selfPerm DbEzeSharePerm
	exits, err := sqlEngine.
		Where("perm_type = ? and token_type = ?", PermTypeAlways, TokenTypeFromSelf).
		Get(&selfPerm)
	if err != nil {
		log.E("get token error")
		panic(err)
	}

	return exits, selfPerm
}

func SaveOtherPerm(token string, deviceId string, deviceName string, deviceType DeviceType) {
	perm := DbEzeSharePerm{
		DeviceName: deviceName,
		DeviceType: deviceType,
		DeviceId:   deviceId,
		Token:      token,
		PermType:   PermTypeAlways,
		TokenType:  TokenTypeFromOther,
	}
	exit, otherPermExits := GetOtherPerm(deviceId)
	if exit {
		perm.Id = otherPermExits.Id
		sqlEngine.ID(perm.Id).
			Cols("device_name", "device_type", "token").
			Update(&otherPermExits)
		// 更新
	} else {
		sqlEngine.InsertOne(&perm)
	}
}

func GetOtherPerm(deviceId string) (bool, DbEzeSharePerm) {
	var otherPerm DbEzeSharePerm
	exits, err := sqlEngine.
		Where("device_id = ?", deviceId).
		Get(&otherPerm)
	if err != nil {
		log.E("get other perm error")
		panic(err)
	}

	return exits, otherPerm
}
