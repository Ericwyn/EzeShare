package storage

import (
	"fmt"
	"github.com/Ericwyn/EzeShare/log"
	"os"
	"xorm.io/xorm"

	_ "github.com/mattn/go-sqlite3"
)

var sqlEngine *xorm.Engine

var dbFileName = "EzeShare.db"

var hadInitDb = false

func InitDb(showSql bool) {
	if hadInitDb {
		return
	}

	var err error

	sqlEngine, err = xorm.NewEngine("sqlite3", GetConfigDirPath()+"/"+dbFileName)
	if err != nil {
		log.E(err)
		log.E("\n\n SQL ENGINE INIT FAIL!!")
		os.Exit(-1)
	}

	// 开启 SQL 打印
	if showSql {
		sqlEngine.ShowSQL(true)
	}

	// 同步表结构
	err = sqlEngine.Sync2(new(DbEzeShareConfig), new(DbEzeShareConnect))
	if err != nil {
		fmt.Println(err)
		log.E("SYNC TABLE ERROR!!")
		os.Exit(-1)
	}

	hadInitDb = true
}
