package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/EzeShare/api/apidef"
	"github.com/Ericwyn/EzeShare/auth"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/EzeShare/utils/deviceutils"
	"github.com/Ericwyn/EzeShare/utils/fileslice"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/GoTools/goSet"
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
	uploadProcessNow := 0

	// 只有文件大小大于最大分片大小的时候，才回去做分片传输，否则也还是用原本的传输方式
	if param.isSlice && param.file.Size() > maxTransferSize {
		// 分块读取实现
		sliceMsgArr := fileslice.Slice(param.file, maxTransferSize)

		successSliceIdSet := goSet.NewSet()

		for i, sliceMsg := range sliceMsgArr {
			if successSliceIdSet.ContinueValue(sliceMsg.SliceNow) {
				log.D("slices ", i, " had transfer success, skip")
				continue
			}

			sliceDataBytes := sliceMsg.ReadSliceToBytes()

			param.sliceMsg = *sliceMsg

			sliceFileName := param.file.Name() + " (slice " + strconv.Itoa(i+1) + "/" + strconv.Itoa(len(sliceMsgArr)) + ")"

			uploadFile := &UploadFile{
				Reader:   bytes.NewReader(*sliceDataBytes),
				FileName: param.file.Name() + "_slice_" + fmt.Sprint(sliceMsg.SliceNow),
				Total:    sliceMsg.SliceSizeBytes,
				TransferPercentCb: func(fileName string, per int) {
					per = calSliceTransferPer(len(sliceMsgArr), i, per)
					if per == uploadProcessNow {
						return
					}
					uploadProcessNow = per
					// 分块计算的传输百分比不太一样
					param.uploadPercentCb(
						sliceFileName,
						uploadProcessNow,
					)
				},
			}
			// 上传这一块的数据
			// 分块上传，对面会在 pubResp 的 data 里面告诉我们有哪些块是已经上传完毕了的
			transferResp := doFileTransferOnce(param, uploadFile)

			var resp apidef.PubResp[[]int64]
			err := json.Unmarshal([]byte(transferResp), &resp)
			if err == nil && resp.Data != nil {
				for _, successSliceId := range resp.Data {
					successSliceIdSet.Put(successSliceId)
				}
			}
		}

	} else {
		param.isSlice = false

		openFile, err := os.Open(param.file.AbsPath())
		if err != nil {
			log.E("open file error", err)
			return
		}

		uploadFile := &UploadFile{
			Reader:   openFile,
			FileName: param.file.Name(),
			Total:    param.file.Size(),
			TransferPercentCb: func(fileName string, per int) {
				if per > uploadProcessNow {
					param.uploadPercentCb(fileName, per)
					uploadProcessNow = per
				}
			},
		}
		// 整块上传
		doFileTransferOnce(param, uploadFile)
	}
}

// calSliceTransferPer 计算分片传输的总百分比
// 比如一共有 10 个分片，你传输了第 6 个，第 7 个传输了 50%, 那么总的传输进度应该是 65%
func calSliceTransferPer(totalSliceNum int, sliceNow int, slicePer int) int {
	perNew := (1.0 / float64(totalSliceNum)) * (float64(sliceNow) + (float64(slicePer) / 100))
	return int(perNew * 100)
}

func doFileTransferOnce(param fileTransferReqParam, uploadFile *UploadFile) string {
	url := "http://" + param.ipAddr + ":" + strconv.Itoa(apidef.HttpApiServerPort) +
		apidef.ApiPathFileTransfer

	unixTimeStamp := time.Now().Unix()

	httpBody := &bytes.Buffer{}
	writer := multipart.NewWriter(httpBody)
	part, err := writer.CreateFormFile("file", param.file.Name())
	if err != nil {
		return ""
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
		return ""
	}

	req, err := http.NewRequest("POST", url, httpBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.E("do file transfer http request error")
		log.E(err)
		return ""
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.E("parse file transfer http request resp error")
			log.E(err)
		}
		resp.Body.Close()
		log.D("file transfer resp status: ", resp.StatusCode)
		log.D("file transfer resp header: ", resp.Header)
		log.D("file transfer resp body: \n", body)

		return body.String()
	}
}
