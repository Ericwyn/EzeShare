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
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

// apiclient 给 sender 请求 receiver 的 api 接口的工具

// DoPermRequest 发起一个文件发送请求
func DoPermRequest(receiverMsg scan.BroadcastMsg, file file.File, permType apidef.PermType, uploadPercentCb func(per int)) {
	alwaysToken := auth.CheckReceiverAlwaysToken(receiverMsg.DeviceId)
	if alwaysToken != "" {
		DoFileTransfer(receiverMsg.Address, apidef.PermTypeAlways, alwaysToken, "", file, uploadPercentCb)
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

	//log.D("DoPermRequest response Status:", httpResp.Status)
	//log.D("DoPermRequest response Headers:", httpResp.Header)
	respBody, _ := io.ReadAll(httpResp.Body)
	//fmt.Println("response Body:", string(body))

	var apiResp apidef.PubResp
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		log.E("parse respBody error: ", string(respBody))
		return
	}
	if apiResp.Code != apidef.RespCodeSuccess {
		log.E("perm req error, ", string(respBody))
		return
	}
	log.D("perm resp : ", string(respBody))

	data := apiResp.Data.(map[string]interface{})

	secTokenResp := data["SecToken"].(string)
	permTypeResp := data["PermType"].(string)
	transferIdResp := data["TransferId"].(string)
	receiverDeviceId := data["ReceiverDeviceId"].(string)
	log.D("perm resp, secToken: ", secTokenResp,
		", permType: ", permTypeResp,
		", transferId: ", transferIdResp)

	decryptToken, err := auth.DecryptRSA(secTokenResp, auth.GetRsaPrivateKeyPath())
	if err != nil {
		return
	}
	log.D("decrypt token: ", decryptToken)
	if permTypeResp == string(apidef.PermReqRespAllowAlways) {
		auth.SaveReceiverAlwaysToken(receiverDeviceId, decryptToken, receiverMsg.Name, receiverMsg.DeviceType)
		// 保存 token 进去
	}
	var permTypeReq apidef.PermType
	if permTypeResp == string(apidef.PermReqRespAllowAlways) {
		permTypeReq = apidef.PermTypeAlways
	} else if permTypeResp == string(apidef.PermReqRespAllowOnce) {
		permTypeReq = apidef.PermTypeOnce
	}

	DoFileTransfer(receiverMsg.Address, permTypeReq, decryptToken, transferIdResp, file, uploadPercentCb)
}

type UploadFile struct {
	io.Reader                        // 读取器
	FileName           string        // 文件名字
	Total              int64         // 总大小
	Current            int64         // 当前大小
	TransferPercentNow int           // 传输百分比
	TransferPercentCb  func(per int) // 传输百分比的回调
	TransferFinishCb   func()        // 传输完成的回调
}

// 实现io.Reader接口的Read方法
// p是一个字节切片，n是读取的字节数，err是错误信息
func (f *UploadFile) Read(p []byte) (n int, err error) {
	n, err = f.Reader.Read(p)
	f.Current += int64(n)
	// 这里可以打印下载进度
	percent := float64(f.Current*10000/f.Total) / 100
	if int(percent) > f.TransferPercentNow {
		f.TransferPercentNow = int(percent)
		f.TransferPercentCb(f.TransferPercentNow)
	}
	if f.Current == f.Total {
		log.D("file upload finish!!!!!!")
		if f.TransferFinishCb != nil {
			f.TransferFinishCb()
		}
	}
	return
}

func DoFileTransfer(ipAddr string,
	permTypeReq apidef.PermType,
	decryptToken string,
	transferId string,
	file file.File,
	uploadPercentCb func(per int),
) {
	url := "http://" + ipAddr + ":" + strconv.Itoa(apidef.HttpApiServerPort) +
		apidef.ApiPathFileTransfer

	unixTimeStamp := time.Now().Unix()

	openFile, err := os.Open(file.AbsPath())
	if err != nil {
		return
	}

	httpBody := &bytes.Buffer{}
	writer := multipart.NewWriter(httpBody)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return
	}

	uploadFile := &UploadFile{
		Reader:            openFile,
		Total:             file.Size(),
		TransferPercentCb: uploadPercentCb,
	}

	_, err = io.Copy(part, uploadFile)

	otherParamMap := map[string]string{
		"sign":         auth.FileTransferSign(decryptToken, file.Name(), unixTimeStamp),
		"transferId":   transferId,
		"fileName":     file.Name(),
		"senderName":   deviceutils.GetDeviceName(),
		"fileSizeBits": strconv.FormatInt(file.Size(), 10),
		"permType":     string(permTypeReq),
		"timeStamp":    strconv.Itoa(int(unixTimeStamp)),
	}
	log.D("do file transfer req, params : ", otherParamMap)
	for k, v := range otherParamMap {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, httpBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.E("do file transfer http request error")
		log.E(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.E("parse file transfer http request resp error")
			log.E(err)
		}
		resp.Body.Close()
		log.D(resp.StatusCode)
		log.D(resp.Header)
		log.D(body)
	}
}

func toJson(obj any) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(marshal)
}
