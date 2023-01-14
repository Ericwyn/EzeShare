package apiclient

import (
	"bytes"
	"encoding/json"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/utils/netutils"
	"io"
	"net/http"
)

// apiclient 给 sender 请求 receiver 的 api 接口的工具

// DoPermRequest 发起一个文件发送请求
func DoPermRequest(addr string, fileName string, fileSizeKb int64, permType apidef.PermType) {
	url := addr + apidef.ApiPathPermReq

	reqStruct := apidef.ApiPermReq{
		PermType:     permType,
		FileName:     fileName,
		FileSizeKb:   fileSizeKb,
		SenderName:   netutils.GetDeviceName(),
		SenderPubKey: auth.GetRsaPublicKey(),
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(toJson(reqStruct))))
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
}

func DoFileTransfer() {
	// TODO client 发送请求
}

func toJson(obj any) string {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(marshal)
}
