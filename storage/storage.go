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

// SaveToken 保存 Token
func SaveToken(token string) {
	tokenConfig := DbEzeShareConfig{
		Key:   ConfigKeyToken,
		Value: token,
	}
	_, err := sqlEngine.InsertOne(tokenConfig)
	if err != nil {
		log.E("save token error")
		panic(err)
	}
}

// GetTokenFromDB 获取 token，存在的话为 bool, string
func GetTokenFromDB() (bool, string) {
	var tokenConfig DbEzeShareConfig
	exits, err := sqlEngine.Where("`key` = ?", ConfigKeyToken).Get(&tokenConfig)
	if err != nil {
		log.E("get token error")
		panic(err)
	}

	return exits, tokenConfig.Value
}
