package storage

import "github.com/Ericwyn/EzeShare/log"

// perm 表相关操作

// SaveSelfToken 保存 OnceToken
func SaveSelfToken(token string) {
	tokenSelf := DbEzeSharePerm{
		DeviceName: "Self",
		DeviceType: "Self",
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

// GetSelfTokenFromDB 获取 token，存在的话为 bool, string
func GetSelfTokenFromDB() (bool, string) {
	var tokenSelf DbEzeSharePerm
	exits, err := sqlEngine.
		Where("perm_type = ? and token_type = ?", PermTypeAlways, TokenTypeFromSelf).
		Get(&tokenSelf)
	if err != nil {
		log.E("get token error")
		panic(err)
	}
	return exits, tokenSelf.Token
}
