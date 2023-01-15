package api

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strconv"
)

func apiReceiver(ctx *gin.Context) {
	sign, signExit := ctx.GetPostForm("sign")
	transferId, transferIdExit := ctx.GetPostForm("transferId")
	fileNameParam, fileNameExit := ctx.GetPostForm("fileName")
	timeStampParam, timeStampExit := ctx.GetPostForm("timeStamp")
	permType, permTypeExit := ctx.GetPostForm("permType")

	if !signExit || !transferIdExit || !timeStampExit || !permTypeExit || !fileNameExit {
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "sign or transferId or timeStamp or permType or fileName param is empty",
		})
		return
	}

	timeStampSec, err := strconv.ParseInt(timeStampParam, 10, 64)
	if err != nil {
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "timeStamp param parse error",
		})
		return
	}
	var signCheck = ""
	if permType == string(apidef.PermTypeOnce) {
		// 通过 transferId 查找到这一条 preSend 记录
		transferMsg := storage.GetTransferMsgFromDB(transferId)
		token := transferMsg.OnceToken
		if token == "" {
			ctx.JSON(200, apidef.PubResp{
				Code: apidef.RespCodeParamError,
				Msg:  "sign error",
			})
			log.I("token error")
			return
		}
		signCheck = auth.FileTransferSign(token, fileNameParam, timeStampSec)
	} else if permType == string(apidef.PermTypeAlways) {
		token := auth.GetSelfToken()
		signCheck = auth.FileTransferSign(token, fileNameParam, timeStampSec)
	} else {
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "perm type param error",
		})
		return
	}

	if signCheck != sign || signCheck == "" {
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "sign error",
		})
		log.I("sign error, true sign : ", signCheck, ", param sign: "+sign)
		return
	}

	// 验证通过, 开始读取和保存文件
	uploadFile, err := ctx.FormFile("file")
	if err != nil {
		log.E("read file error")
		log.E(err)
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "read file error",
		})
		return
	}
	saveUploadFile(ctx, uploadFile, transferId)

}

func saveUploadFile(ctx *gin.Context, uploadFile *multipart.FileHeader, transferId string) {
	fileName := uploadFile.Filename
	saveDirPath := storage.GetDownloadDirPath()

	finalSavePath := saveDirPath + "/" + fileName
	finalSaveFile := file.OpenFile(finalSavePath)
	if finalSaveFile.Exits() {
		newFileName := ""
		// 从 1 开始拼接后缀，直到找到为空的路径
		for i := 1; i <= 50; i++ {
			tryPath := finalSavePath + "(" + strconv.Itoa(i) + ")"
			tryPathFile := file.OpenFile(tryPath)
			if !tryPathFile.Exits() {
				newFileName = fileName + "(" + strconv.Itoa(i) + ")"
				finalSavePath = tryPath
				finalSaveFile = tryPathFile
				log.I("save file to new name: ", newFileName)
				// 如果这个路径不存在，证明这个文件名可以用
				break
			}
		}
		if newFileName == "" {
			log.E("can't find the new name of upload file, " +
				"there are to many same name file in the down dir")

			ctx.JSON(200, apidef.PubResp{
				Code: apidef.RespCodeServerError,
				Msg: "can't find the new name of upload file, " +
					"there are to many same name file in the down dir",
			})
			return
		}
		// 文件保存记录得更新
		storage.RenameUploadFileToDB(transferId, newFileName)
	}
	// 开始传输
	storage.SaveTransferStatus(transferId, storage.TransferStatusSending, finalSavePath, uploadFile.Size/1024)
	err := ctx.SaveUploadedFile(uploadFile, finalSavePath)
	storage.SaveTransferStatus(transferId, storage.TransferStatusFinish, "", 0)
	if err != nil {
		log.E("ctx save upload file fail")
		log.E(err)
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeServerError,
			Msg:  "ctx save upload file fail",
		})
		return
	}
	log.I("save file success, filePath : ", finalSavePath)
	ctx.JSON(200, apidef.PubResp{
		Code: apidef.RespCodeSuccess,
		Msg:  "transfer file success",
	})
	return
}
