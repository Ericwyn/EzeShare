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
