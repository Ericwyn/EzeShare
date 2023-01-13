package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/EzeShare/api"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/utils/netutils"
	"io"
	"net/http"
)

// apiclient 给 sender 请求 receiver 的 api 接口的工具

// DoPermRequest 发起一个文件发送请求
func DoPermRequest(addr string, fileName string, fileSizeKb int64, permType api.PermType) {
	url := addr + api.ApiPathPermReq

	reqStruct := api.ApiPermReq{
		PermType:     permType,
		FileName:     fileName,
		FileSizeKb:   0,
		SenderName:   netutils.GetDeviceName(),
		SenderPubKey: auth.GetRsaPublicKey(),
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(toJson(reqStruct))))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.D("DoPermRequest response Status:", resp.Status)
	log.D("DoPermRequest response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	// TODO client 发送请求
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
