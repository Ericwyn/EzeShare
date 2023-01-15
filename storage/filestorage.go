package storage

import (
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/GoTools/file"
	"os"
)

// 用来实现本地 IO 相关方法

const configDirName = ".EzeShare"
const fileSaveDirName = "Downloads/EzeShareFiles"

var configDirPathCache = ""
var fileSaveDirNameCache = ""

// GetConfigDirPath 获取配置文件目录， 目录不存在的话就会创建
func GetConfigDirPath() string {
	if configDirPathCache != "" {
		return configDirPathCache
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
	configDirPathCache = configDir.AbsPath()
	log.D("get config dir path: " + configDirPathCache)
	return configDirPathCache
}

func GetDownloadDirPath() string {
	if fileSaveDirNameCache != "" {
		return fileSaveDirNameCache
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.E("get config homeDir error, can' get user home homeDir")
		panic(err)
	}

	downloadDir := file.OpenFile(homeDir + "/" + fileSaveDirName)
	if !downloadDir.Exits() {
		downloadDir.Mkdirs()
	}
	fileSaveDirNameCache = downloadDir.AbsPath()
	log.D("get config dir path: " + fileSaveDirNameCache)
	return fileSaveDirNameCache
}
