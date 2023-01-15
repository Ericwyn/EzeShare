package storage

import "github.com/Ericwyn/EzeShare/log"

// perm 表相关操作

// SaveSelfToken 保存 OnceToken
func SaveSelfToken(token string, deviceId string) {
	tokenSelf := DbEzeSharePerm{
		DeviceName: "Self",
		DeviceType: "Self",
		DeviceID:   deviceId,
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
