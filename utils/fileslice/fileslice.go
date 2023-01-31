package fileslice

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/Ericwyn/EzeShare/log"
	"github.com/Ericwyn/GoTools/file"
)

type SliceMsg struct {
	TotalSliceNum     int64  `json:"totalSliceNum"`     // 分块数量
	SliceNow          int64  `json:"sliceNow"`          // 当前分块 id
	MaxSliceSizeBytes int64  `json:"maxSliceSizeBytes"` // 最大分块大小
	SliceSizeBytes    int64  `json:"sliceSizeBytes"`    // 当前分块大小
	FileSizeBytes     int64  `json:"fileSizeBytes"`     // 所有分块大小 (文件大小)
	FileName          string `json:"fileName"`          // 文件名称

	fileMsg     file.File // 文件信息
	offsetBytes int64     // 开始读取的 bytes
	endBytes    int64     // 结束读取的 bytes
	md5Sum      string    // 当前分块 md5 信息
	fileData    *[]byte   // 数据
}

// ReadSliceToBytes 读取当前分块信息
func (sliceMsg *SliceMsg) ReadSliceToBytes() *[]byte {
	if sliceMsg.fileData != nil && sliceMsg.md5Sum != "" {
		return sliceMsg.fileData
	}
	// 读取 data, 并且计算 md5
	openFile, err := sliceMsg.fileMsg.Open()
	if err != nil {
		return nil
	}

	readByte := make([]byte, sliceMsg.SliceSizeBytes)

	_, err = openFile.ReadAt(readByte, sliceMsg.offsetBytes)
	if err != nil {
		log.E("slice file error, all slice num: ", sliceMsg.TotalSliceNum,
			", total slice: ", sliceMsg.SliceNow)
		panic(err)
		return nil
	}

	//log.D("file slice read success, target slice bytes: ", sliceMsg.SliceSizeBytes,
	//	", read slice: ", readLen)

	sliceMsg.fileData = &readByte

	// 计算 md5
	sliceMsg.md5Sum = calMd5(*sliceMsg.fileData)

	return &readByte
}

func (sliceMsg *SliceMsg) GetSliceMd5Sum() string {
	if sliceMsg.md5Sum != "" {
		return sliceMsg.md5Sum
	}
	// 需要读取文件，并且计算 MD5
	sliceMsg.ReadSliceToBytes()
	return sliceMsg.md5Sum
}

// Slice 将一个文件拆分成多个
func Slice(fileMsg file.File, maxSliceSizeBits int64) []*SliceMsg {
	fileSizeBytes := fileMsg.Size()
	totalSliceNum := fileSizeBytes / maxSliceSizeBits // 总分块数量
	if fileSizeBytes%maxSliceSizeBits != 0 {
		totalSliceNum += 1
	}

	// 第几字节
	readOffset := int64(0)

	resArr := make([]*SliceMsg, totalSliceNum)

	for i := int64(0); i < totalSliceNum; i++ {
		endBytes := readOffset + maxSliceSizeBits
		if endBytes > fileSizeBytes {
			endBytes = fileSizeBytes
		}

		msg := SliceMsg{
			TotalSliceNum:     totalSliceNum,
			SliceNow:          i,
			MaxSliceSizeBytes: maxSliceSizeBits,
			SliceSizeBytes:    endBytes - readOffset,
			FileSizeBytes:     fileSizeBytes,
			FileName:          fileMsg.Name(),

			fileMsg:     fileMsg,
			offsetBytes: readOffset,
			endBytes:    endBytes,
		}

		resArr[i] = &msg

		readOffset += maxSliceSizeBits
	}

	return resArr
}

func calMd5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
