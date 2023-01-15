package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/utils/netutils"
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
func DoPermRequest(ipAddr string, file file.File, permType apidef.PermType) {
	url := "http://" + ipAddr + ":" + strconv.Itoa(apidef.HttpApiServerPort) +
		apidef.ApiPathPermReq

	reqStruct := apidef.ApiPermReq{
		PermType:     permType,
		FileName:     file.Name(),
		FileSizeBits: file.Size(),
		SenderName:   netutils.GetDeviceName(),
		SenderPubKey: auth.GetRsaPublicKey(),
	}

	reqJson := toJson(reqStruct)

	log.I("request file transfer perm to ", url)
	log.I("request body : ", reqJson)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqJson)))
	if err != nil {
		log.I("new http req error")
		log.I(err)
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
		log.I("parse respBody error: ", string(respBody))
		return
	}
	if apiResp.Code != apidef.RespCodeSuccess {
		log.E("perm req error, ", string(respBody))
		return
	}
	data := apiResp.Data.(map[string]interface{})

	secTokenResp := data["SecToken"].(string)
	permTypeResp := data["PermType"].(string)
	transferIdResp := data["TransferId"].(string)
	log.I("perm resp, secToken: ", secTokenResp,
		", permType: ", permTypeResp,
		", transferId: ", transferIdResp)

	decryptToken, err := auth.DecryptRSA(secTokenResp, auth.GetRsaPrivateKeyPath())
	if err != nil {
		return
	}
	log.I("decrypt token: ", decryptToken)
	DoFileTransfer(ipAddr, decryptToken, transferIdResp, file)
}

func DoFileTransfer(ipAddr string, decryptToken string, transferId string, file file.File) {
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
	_, err = io.Copy(part, openFile)

	_ = writer.WriteField("sign", auth.FileTransferSign(decryptToken, file.Name(), unixTimeStamp))
	_ = writer.WriteField("transferId", transferId)
	_ = writer.WriteField("timeStamp", strconv.Itoa(int(time.Now().Unix())))

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
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}

func toJson(obj any) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(marshal)
}
