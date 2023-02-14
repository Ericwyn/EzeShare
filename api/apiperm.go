package api

import (
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/conf"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/storage"
	"github.com/Ericwyn/EzeShare/utils/errutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go/types"
	"sync"
	"time"
)

// 一次只允许一个 perm check 请求
var isPermCheckNow = false

type PermReqCallback func(apidef.ApiPermReq) apidef.PermReqRespType

var permReqCallback PermReqCallback

func SetPermReqCallback(cb PermReqCallback) {
	permReqCallback = cb
}

func apiPermReq(ctx *gin.Context) {
	if isPermCheckNow {
		ctx.JSON(200, apidef.PubResp[types.Nil]{
			Code: apidef.RespCodeServerError,
			Msg:  "server check another perm req now, please try again later",
		})
		return
	}

	isPermCheckNow = true
	defer func() { isPermCheckNow = false }()

	var reqBody apidef.ApiPermReq
	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.JSON(200, apidef.PubResp[types.Nil]{
			Code: apidef.RespCodeParamError,
			Msg:  "json parse error",
		})
		return
	}

	// 校验参数
	if checkReq := reqBody.CheckReq(); checkReq != "" {
		ctx.JSON(200, apidef.PubResp[types.Nil]{
			Code: apidef.RespCodeParamError,
			Msg:  checkReq,
		})
		return
	}

	// 设置 req addr
	reqBody.SenderAddr = ctx.ClientIP()

	if conf.RunInPermCheckMode {
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
				log.E("wait group done panic, ", i)
			})
			permRespTypeFromUI = permType
		}()

		// 开始阻塞等待回调
		log.D("start wait for perm req ui callback")
		wg.Wait()

		// 阻塞结束进行处理
		log.D("perm req callback wait done, result: ", permRespTypeFromUI)

		ctx.JSON(200, generalPermResp(reqBody, permRespTypeFromUI))
	} else {
		log.I("自动允许收到来自", reqBody.SenderName, "的文件: ", reqBody.FileName,
			", 大小: ", reqBody.FileSizeBytes)
		ctx.JSON(200, generalPermResp(reqBody, apidef.PermReqRespAllowOnce))
	}

}

func generalPermResp(reqBody apidef.ApiPermReq, permRespType apidef.PermReqRespType) any {
	if permRespType == apidef.PermReqRespAllowAlways {
		// 拿到自己的 token
		alwaysToken := auth.GetSelfToken()
		// 公钥加密
		secToken, err := auth.EncryptRSAWithPubKeyStr(alwaysToken, reqBody.SenderPubKey)
		if err != nil {
			return apidef.PubResp[types.Nil]{
				Code: apidef.RespCodeParamError,
				Msg:  "encrypt with rsa pub key error",
			}
		}

		transferId := uuid.New().String()
		// 写一条记录进去数据库
		transferMsg := storage.DbEzeShareTransferMsg{
			TransferId:        transferId,
			FileName:          reqBody.FileName,
			FileSizeKb:        reqBody.FileSizeBytes,
			OnceToken:         "",
			TransferStatus:    storage.TransferStatusPreSend,
			FromDeviceName:    reqBody.SenderName,
			FromDeviceAddress: reqBody.SenderAddr,
			RequestTime:       time.Now(),
		}
		storage.SavePreTransferMsg(transferMsg)

		return apidef.PubResp[apidef.ApiPermResp]{
			Code: apidef.RespCodeSuccess,
			Msg:  "success",
			Data: apidef.ApiPermResp{
				SecToken:         secToken,
				PermType:         permRespType,
				TransferId:       transferId,
				ReceiverDeviceId: auth.GetSelfDeviceId(),
			},
		}
	} else if permRespType == apidef.PermReqRespAllowOnce {
		onceToken := uuid.New().String()
		transferId := uuid.New().String()

		// 公钥加密
		secToken, err := auth.EncryptRSAWithPubKeyStr(onceToken, reqBody.SenderPubKey)
		if err != nil {
			return apidef.PubResp[types.Nil]{
				Code: apidef.RespCodeParamError,
				Msg:  "encrypt with rsa pub key error",
			}
		}

		// 写一条记录进去数据库
		transferMsg := storage.DbEzeShareTransferMsg{
			TransferId:        transferId,
			FileName:          reqBody.FileName,
			FileSizeKb:        reqBody.FileSizeBytes,
			OnceToken:         onceToken,
			TransferStatus:    storage.TransferStatusPreSend,
			FromDeviceName:    reqBody.SenderName,
			FromDeviceAddress: reqBody.SenderAddr,
			RequestTime:       time.Now(),
		}
		storage.SavePreTransferMsg(transferMsg)

		return apidef.PubResp[apidef.ApiPermResp]{
			Code: apidef.RespCodeSuccess,
			Msg:  "success",
			Data: apidef.ApiPermResp{
				SecToken:         secToken,
				TransferId:       transferId,
				PermType:         permRespType,
				ReceiverDeviceId: auth.GetSelfDeviceId(),
			},
		}

	} else {
		return apidef.PubResp[apidef.ApiPermResp]{
			Code: apidef.RespCodeParamError,
			Msg:  "",
			Data: apidef.ApiPermResp{
				SecToken:         "",
				PermType:         apidef.PermReqRespDisAllow,
				ReceiverDeviceId: auth.GetSelfDeviceId(),
			},
		}
	}
}
