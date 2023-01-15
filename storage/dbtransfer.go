package storage

import (
	"github.com/Ericwyn/EzeShare/log"
	"time"
)

// transfer msg 表相关操作

// SavePreTransferMsg 保存一条传输记录, 签发 OnceToken
func SavePreTransferMsg(transMsg DbEzeShareTransferMsg) {
	// 设置一些初始数据
	_, err := sqlEngine.InsertOne(transMsg)
	if err != nil {
		log.E("save pre transfer msg error")
		return
	}
	log.D("save transfer msg success, id: ", transMsg.Id)
}

// SaveTransferStatus 更改传输的状态
func SaveTransferStatus(transferId string, status TransferStatus, fileSavePath string, fileSizeKb int64) {
	columns := []string{"transfer_status"}
	transMsg := DbEzeShareTransferMsg{
		TransferId:     transferId,
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

	if fileSavePath != "" {
		transMsg.FileSavePath = fileSavePath
		columns = append(columns, "file_save_path")
	}
	if fileSizeKb != 0 {
		transMsg.FileSizeKb = fileSizeKb
		columns = append(columns, "file_size_kb")
	}

	affected, err := sqlEngine.Where("transfer_id = ?", transferId).
		Cols(columns...).
		Update(&transMsg)
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

func GetTransferMsgFromDB(transferId string) *DbEzeShareTransferMsg {
	var transferMsg DbEzeShareTransferMsg
	exits, err := sqlEngine.Where("transfer_id = ?", transferId).Get(&transferMsg)
	if err != nil {
		log.E("get db transfer msg error")
		log.E(err)
		return nil
	}
	if !exits {
		log.E("can't find transfer msg by transferId : ", transferMsg)
		return nil
	}
	return &transferMsg
}

// RenameUploadFileToDB 上传的文件遇到同名的情况, 需要更新名字
func RenameUploadFileToDB(transferId string, newFileName string) {
	transferMsg := DbEzeShareTransferMsg{
		FileName: newFileName,
	}
	affected, err := sqlEngine.Where("transfer_id = ?", transferId).
		Cols("file_name").
		Update(&transferMsg)
	if err != nil {
		log.E("update file name error, transfer_id: ", transferId, ", new file name: ", newFileName)
		return
	}
	if affected != 1 {
		log.E("update file name file, transfer_id: ", transferId, ", new file name: ", newFileName)
		return
	}
}
