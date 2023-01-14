package storage

import (
	"github.com/Ericwyn/EzeShare/log"
	"time"
)

// transfer msg 表相关操作

// SavePreTransferMsg 保存一条传输记录, 签发 OnceToken
func SavePreTransferMsg(transMsg DbEzeShareTransferMsg) {
	// 设置一些初始数据
	id, err := sqlEngine.InsertOne(transMsg)
	if err != nil {
		log.E("save pre transfer msg error")
		return
	}
	log.I("save transfer msg success, id: ", id)
}

// SaveTransferStatus 更改传输的状态
func SaveTransferStatus(transMsgId int64, status TransferStatus) {
	columns := []string{"transfer_status"}
	transMsg := DbEzeShareTransferMsg{
		Id:             transMsgId,
		TransferStatus: status,
	}

	//if status == TransferStatusPreSend {
	//	transMsg.StartTime = time.Now()
	//	columns = append(columns, "request_time")
	//} else
	if status == TransferStatusSending {
		transMsg.StartTime = time.Now()
		columns = append(columns, "start_time")
	} else if status == TransferStatusFinish {
		transMsg.FinishTime = time.Now()
		columns = append(columns, "finish_time")
	}

	affected, err := sqlEngine.ID(transMsgId).Cols(columns...).Update(&transMsg)
	if err != nil {
		log.E("update transfer status error, id: ", transMsg, "status: ", status)
		return
	}
	if affected != 1 {
		log.E("update transfer status fail, id: ", transMsg, "status: ", status)
		return
	}
	log.I("update transfer status, id: ", transMsg, "status: ", status)
}
