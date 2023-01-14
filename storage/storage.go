package storage

import (
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/GoTools/file"
	"os"
)

const configDirName = ".EzeShare"

var configDirPath = ""

// GetConfigDirPath 获取配置文件目录， 目录不存在的话就会创建
func GetConfigDirPath() string {
	if configDirPath != "" {
		return configDirPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.E("get config homeDir error, can' get user home homeDir")
		panic(err)
	}

	configDir := file.OpenFile(homeDir + "/" + configDirName)
	if !configDir.Exits() {
		configDir.Mkdirs()
	}
	configDirPath = configDir.AbsPath()
	log.D("get config dir path: " + configDirPath)
	return configDirPath
}

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
