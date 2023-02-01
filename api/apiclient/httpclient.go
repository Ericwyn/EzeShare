package apiclient

import (
	"bytes"
	"encoding/json"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/scan"
	"github.com/Ericwyn/EzeShare/utils/deviceutils"
	"github.com/Ericwyn/GoTools/file"
	"io"
	"net/http"
	"strconv"
)

// apiclient 给 sender 请求 receiver 的 api 接口的工具

// DoPermRequest 发起一个文件发送请求
func DoPermRequest(receiverMsg scan.BroadcastMsg,
	file file.File,
	permType apidef.PermType,
	uploadPercentCb func(fileName string, per int),
) {
	alwaysToken := auth.CheckReceiverAlwaysToken(receiverMsg.DeviceId)
	if alwaysToken != "" {
		parm := fileTransferReqParam{
			ipAddr:          receiverMsg.Address,
			permTypeReq:     apidef.PermTypeAlways,
			decryptToken:    alwaysToken,
			transferId:      "",
			file:            file,
			uploadPercentCb: uploadPercentCb,
		}
		DoFileTransfer(parm)
		return
	}

	url := "http://" + receiverMsg.Address + ":" + strconv.Itoa(apidef.HttpApiServerPort) +
		apidef.ApiPathPermReq

	reqStruct := apidef.ApiPermReq{
		PermType:     permType,
		FileName:     file.Name(),
		FileSizeBits: file.Size(),
		SenderName:   deviceutils.GetDeviceName(),
		SenderPubKey: auth.GetRsaPublicKey(),
	}

	reqJson := toJson(reqStruct)

	log.D("request file transfer perm to ", url)
	log.D("request body : ", reqJson)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqJson)))
	if err != nil {
		log.E("new http req error")
		log.E(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		log.E("do http request error")
		log.E(err)
		return
	}
	defer httpResp.Body.Close()

	respBody, _ := io.ReadAll(httpResp.Body)

	var apiPubResp apidef.PubResp[apidef.ApiPermResp]
	err = json.Unmarshal(respBody, &apiPubResp)
	if err != nil {
		log.E("parse respBody error: ", string(respBody))
		return
	}
	if apiPubResp.Code != apidef.RespCodeSuccess {
		log.E("perm req error, ", string(respBody))
		return
	}
	log.D("perm resp : ", string(respBody))

	//data := apiPubResp.Data
	//
	//respJson := toJson(data)
	var apiPermResp apidef.ApiPermResp = apiPubResp.Data
	//err = json.Unmarshal([]byte(respJson), &apiPermResp)
	//if err != nil {
	//	log.E("parse api perm resp json error")
	//	return
	//}

	log.D("sec token: " + apiPermResp.SecToken)
	decryptToken, err := auth.DecryptRSA(apiPermResp.SecToken, auth.GetRsaPrivateKeyPath())
	if err != nil {
		log.E("decrypt sec token fail", err)
		return
	}
	log.D("decrypt token: ", decryptToken)
	if apiPermResp.PermType == apidef.PermReqRespAllowAlways {
		auth.SaveReceiverAlwaysToken(apiPermResp.ReceiverDeviceId, decryptToken, receiverMsg.Name, receiverMsg.DeviceType)
		// 保存 token 进去
	}
	var permTypeReq apidef.PermType
	//if apiPermResp.PermType == apidef.PermReqRespAllowAlways {
	//	permTypeReq = apidef.PermTypeAlways
	//} else if apiPermResp.PermType == apidef.PermReqRespAllowOnce {
	//	permTypeReq = apidef.PermTypeOnce
	//} else {
	//	log.D("permType error: " + apiPermResp.PermType)
	//}

	if apiPermResp.PermType.Equals(apidef.PermReqRespAllowAlways) {
		permTypeReq = apidef.PermTypeAlways
	} else if apiPermResp.PermType.Equals(apidef.PermReqRespAllowOnce) {
		permTypeReq = apidef.PermTypeOnce
	} else {
		log.D("permType error: " + apiPermResp.PermType)
	}
	parm := fileTransferReqParam{
		ipAddr:          receiverMsg.Address,
		permTypeReq:     permTypeReq,
		decryptToken:    decryptToken,
		transferId:      apiPermResp.TransferId,
		file:            file,
		uploadPercentCb: uploadPercentCb,
		isSlice:         receiverMsg.DeviceType == "Android",
	}
	DoFileTransfer(parm)
}

type UploadFile struct {
	io.Reader                                        // 读取器
	Total             int64                          // 总大小
	TransferPercentCb func(fileName string, per int) // 传输百分比的回调
	TransferFinishCb  func()                         // 传输完成的回调

	FileName           string // 文件名字
	currentReadBytes   int64  // 当前大小
	transferPercentNow int    // 传输百分比
}

// 实现io.Reader接口的Read方法
// p是一个字节切片，n是读取的字节数，err是错误信息
func (f *UploadFile) Read(p []byte) (n int, err error) {
	n, err = f.Reader.Read(p)
	f.currentReadBytes += int64(n)
	// 这里可以打印下载进度
	percent := float64(f.currentReadBytes*10000/f.Total) / 100
	if int(percent) > f.transferPercentNow {
		f.transferPercentNow = int(percent)
		f.TransferPercentCb(f.FileName, f.transferPercentNow)
	}
	if f.currentReadBytes == f.Total {
		log.D("file upload finish!!!!!!")
		if f.TransferFinishCb != nil {
			f.TransferFinishCb()
		}
	}
	return
}

func toJson(obj any) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(marshal)
}
