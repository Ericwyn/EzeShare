package api

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/EzeShare/utils/errutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sync"
	"time"
)

type PermReqCallback func(apidef.ApiPermReq) apidef.PermReqRespType

var permReqCallback PermReqCallback

func SetPermReqCallback(cb PermReqCallback) {
	permReqCallback = cb
}

func apiPermReq(ctx *gin.Context) {
	var reqBody apidef.ApiPermReq
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "json parse error",
		})
		return
	}

	// 校验参数
	if checkReq := reqBody.CheckReq(); checkReq != "" {
		ctx.JSON(200, apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  checkReq,
		})
		return
	}

	// 设置 req addr
	reqBody.SenderAddr = ctx.ClientIP()

	wg := sync.WaitGroup{}
	wg.Add(1)

	wgDone := false
	var permRespTypeFromUI apidef.PermReqRespType

	// 60s 后解锁
	go func() {
		time.Sleep(60 * time.Second)
		errutils.Try(func() {
			if !wgDone {
				wg.Done()
				wgDone = true
			}
			permRespTypeFromUI = apidef.PermReqRespDisAllow
		}, func(i interface{}) {
			log.I("wait group done panic, ", i)
		})
	}()

	go func() {
		// 等待 ui 选择结果，此处阻塞
		permType := permReqCallback(reqBody)

		errutils.Try(func() {
			if !wgDone {
				wg.Done()
				wgDone = true
			}
		}, func(i interface{}) {
			log.I("wait group done panic, ", i)
		})
		permRespTypeFromUI = permType
	}()

	// 开始阻塞等待回调
	log.D("start wait for perm req ui callback")
	wg.Wait()

	// 阻塞结束进行处理
	log.D("perm req callback wait done, result: ", permRespTypeFromUI)

	ctx.JSON(200, generalPermResp(reqBody, permRespTypeFromUI))
}

func generalPermResp(reqBody apidef.ApiPermReq, permRespType apidef.PermReqRespType) any {
	if permRespType == apidef.PermReqRespAllowAlways {
		// 拿到自己的 token
		alwaysToken := auth.GetSelfToken()
		// 公钥加密
		secToken, err := auth.EncryptRSAWithPubKeyStr(alwaysToken, reqBody.SenderPubKey)
		if err != nil {
			return apidef.PubResp{
				Code: apidef.RespCodeParamError,
				Msg:  "encrypt with rsa pub key error",
			}
		}

		transferId := uuid.New().String()
		// 写一条记录进去数据库
		transferMsg := storage.DbEzeShareTransferMsg{
			TransferId:        transferId,
			FileName:          reqBody.FileName,
			FileSizeKb:        reqBody.FileSizeBits,
			OnceToken:         "",
			TransferStatus:    storage.TransferStatusPreSend,
			FromDeviceName:    reqBody.SenderName,
			FromDeviceAddress: reqBody.SenderAddr,
			RequestTime:       time.Now(),
		}
		storage.SavePreTransferMsg(transferMsg)

		return apidef.PubResp{
			Code: apidef.RespCodeSuccess,
			Msg:  "success",
			Data: apidef.ApiPermResp{
				SecToken:   secToken,
				PermType:   permRespType,
				TransferId: transferId,
			},
		}
	} else if permRespType == apidef.PermReqRespAllowOnce {
		onceToken := uuid.New().String()
		transferId := uuid.New().String()

		// 公钥加密
		secToken, err := auth.EncryptRSAWithPubKeyStr(onceToken, reqBody.SenderPubKey)
		if err != nil {
			return apidef.PubResp{
				Code: apidef.RespCodeParamError,
				Msg:  "encrypt with rsa pub key error",
			}
		}

		// 写一条记录进去数据库
		transferMsg := storage.DbEzeShareTransferMsg{
			TransferId:        transferId,
			FileName:          reqBody.FileName,
			FileSizeKb:        reqBody.FileSizeBits,
			OnceToken:         onceToken,
			TransferStatus:    storage.TransferStatusPreSend,
			FromDeviceName:    reqBody.SenderName,
			FromDeviceAddress: reqBody.SenderAddr,
			RequestTime:       time.Now(),
		}
		storage.SavePreTransferMsg(transferMsg)

		return apidef.PubResp{
			Code: apidef.RespCodeSuccess,
			Msg:  "success",
			Data: apidef.ApiPermResp{
				SecToken:   secToken,
				TransferId: transferId,
				PermType:   permRespType,
			},
		}

	} else {
		return apidef.PubResp{
			Code: apidef.RespCodeParamError,
			Msg:  "",
			Data: apidef.ApiPermResp{
				SecToken: "",
				PermType: apidef.PermReqRespDisAllow,
			},
		}
	}
}
