package fileslice

import (
	"github.com/Ericwyn/GoTools/file"
)

type FileSliceMsg struct {
	TotalSliceNum    int       // 分块数量
	SliceNow         int       // 当前分块 id
	MaxSliceSizeBits int64     // 最大分块大小
	SliceSizeBits    int64     // 当前分块大小
	TotalSizeBits    int64     // 所有分块大小 (文件大小)
	Md5SumList       []string  // 各个分块的 md5sum
	FileMsg          file.File // 文件信息
}

// ReadSliceToBytes 读取当前分块信息
func (sliceMsg *FileSliceMsg) ReadSliceToBytes() []byte {
	return nil
}

// Slice 将一个文件拆分成多个
func Slice(fileMsg file.File, maxSliceSizeBits int64) []FileSliceMsg {
	size := fileMsg.Size()
	totalSliceNum := size / maxSliceSizeBits // 总分块数量
	if size%maxSliceSizeBits != 0 {
		totalSliceNum += 1
	}

	openFile, err := fileMsg.Open()
	if err != nil {
		return nil
	}
	readByte := make([]byte, maxSliceSizeBits)
	readOffset := int64(0)
	//var i int64
	// 先计算出来一堆 md5 先
	for i := int64(0); i < totalSliceNum; i++ {
		// 读取出来这一块，并且计算 md5
		openFile.ReadAt(readByte, readOffset)
		// TODO 分块读取

	}
	return nil
}
