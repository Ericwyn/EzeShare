package fileslice

import (
	"fmt"
	"github.com/Ericwyn/GoTools/file"
	"testing"
)

const _1MB_BYTES int64 = 1048576
const maxTransferSize = 10 * _1MB_BYTES

func TestSlice(t *testing.T) {

	fileMsg := file.OpenFile("C:/Users/Ericwyn/Downloads/PowerToysSetup-0.66.0-x64.exe")

	fmt.Println("fileName: ", fileMsg.Name(), ", size: ", fileMsg.Size())

	// 全部读取出来
	readBytes, err := fileMsg.Read()
	if err != nil {
		panic(err)
	}

	fileSize := int64(len(readBytes))

	sliceIdx := 0
	for i := int64(0); i < fileSize; i += maxTransferSize {
		endIdx := i + maxTransferSize
		if endIdx > fileSize {
			endIdx = fileSize
		}
		byteSlices := readBytes[i:endIdx]
		md5Sum := calMd5(byteSlices)
		fmt.Println("slice:", sliceIdx, ", size: ", len(byteSlices), ", md5: ", md5Sum)

		sliceIdx++
	}

	fmt.Println("-------------------------------------------------------")

	sliceArr := Slice(fileMsg, maxTransferSize)
	for i, fileSliceMsg := range sliceArr {
		fmt.Println("slice:", i, ", size: ", fileSliceMsg.SliceSizeBytes, ", md5: ", fileSliceMsg.GetSliceMd5Sum())
	}
}
