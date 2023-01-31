package apiclient

import (
	"bytes"
	"fmt"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/utils/deviceutils"
	"github.com/Ericwyn/EzeShare/utils/fileslice"
	"github.com/Ericwyn/GoTools/file"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

const _1MB_BYTES int64 = 1048576
const maxTransferSize = 30 * _1MB_BYTES

type fileTransferReqParam struct {
	ipAddr          string
	permTypeReq     apidef.PermType
	decryptToken    string
	transferId      string
	file            file.File
	uploadPercentCb func(fileName string, per int)
	isSlice         bool // 是否分块上传，针对 Android 设备，避免 android 上 netty http service oom
	sliceMsg        fileslice.SliceMsg
}

// DoFileTransfer 文件上传 (一次性)
func DoFileTransfer(param fileTransferReqParam) {
	if param.isSlice {
		// 分块读取实现
		sliceMsgArr := fileslice.Slice(param.file, maxTransferSize)
		for _, sliceMsg := range sliceMsgArr {

			sliceDataBytes := sliceMsg.ReadSliceToBytes()

			param.sliceMsg = *sliceMsg

			uploadFile := &UploadFile{
				Reader:            bytes.NewReader(*sliceDataBytes),
				FileName:          param.file.Name() + "_slice_" + fmt.Sprint(sliceMsg.SliceNow),
				Total:             sliceMsg.SliceSizeBytes,
				TransferPercentCb: param.uploadPercentCb,
			}
			// 整块上传
			doFileTransferOnce(param, uploadFile)
		}

	} else {
		openFile, err := os.Open(param.file.AbsPath())
		if err != nil {
			log.E("open file error", err)
			return
		}
		uploadFile := &UploadFile{
			Reader:            openFile,
			FileName:          param.file.Name(),
			Total:             param.file.Size(),
			TransferPercentCb: param.uploadPercentCb,
		}
		// 整块上传
		doFileTransferOnce(param, uploadFile)
	}
}

func doFileTransferSlice() {

}

func printUploadProcess(fileName string, per int) {
	process := "["
	for i := 0; i <= 100; i += 5 {
		if i < per {
			process += "="
		} else {
			process += " "
		}
	}
	process += "]"
	log.I("上传 ", fileName, ", 进度: ", process)
}

func doFileTransferOnce(param fileTransferReqParam, uploadFile *UploadFile) {
	url := "http://" + param.ipAddr + ":" + strconv.Itoa(apidef.HttpApiServerPort) +
		apidef.ApiPathFileTransfer

	unixTimeStamp := time.Now().Unix()

	httpBody := &bytes.Buffer{}
	writer := multipart.NewWriter(httpBody)
	part, err := writer.CreateFormFile("file", param.file.Name())
	if err != nil {
		return
	}

	//uploadFile := &UploadFile{
	//	Reader:            openFile,
	//	Total:             param.file.Size(),
	//	TransferPercentCb: param.uploadPercentCb,
	//}

	_, err = io.Copy(part, uploadFile)

	otherParamMap := map[string]string{
		"sign":         auth.FileTransferSign(param.decryptToken, param.file.Name(), unixTimeStamp),
		"transferId":   param.transferId,
		"fileName":     param.file.Name(),
		"senderName":   deviceutils.GetDeviceName(),
		"fileSizeBits": strconv.FormatInt(param.file.Size(), 10),
		"permType":     string(param.permTypeReq),
		"timeStamp":    strconv.Itoa(int(unixTimeStamp)),
		"isSlice":      fmt.Sprint(param.isSlice),
	}

	if param.isSlice {
		otherParamMap["sliceMsg"] = toJson(param.sliceMsg)
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
