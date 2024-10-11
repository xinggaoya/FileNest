package impl

import (
	"FileNest/common/model"
	"FileNest/internal/consts"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

/**
  @author: XingGao
  @date: 2024/9/28
**/

type HelloWordServiceImpl struct{}

type SysUser struct {
	model.BaseEntity
	Name string `json:"name" gorm:"type:varchar(50)"`
	Age  int    `json:"age" gorm:"type:int"`
	Sex  string `json:"sex" gorm:"type:varchar(2)"`
	Desc string `json:"desc" gorm:"type:varchar(150)"`
}

func (h *HelloWordServiceImpl) GetFileList(path string) {

}

func (h *HelloWordServiceImpl) UploadFile(file *multipart.FileHeader, chunkIndex, totalChunks int) error {
	// Create the file to store the chunk
	chunkPath := filepath.Join(consts.TempDir, "chunk_"+strconv.Itoa(chunkIndex))
	outFile, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("Unable to create chunk file", err)
	}
	defer outFile.Close()

	// Write the chunk to the file
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("Unable to open chunk file", err)
	}
	defer f.Close()
	_, err = io.Copy(outFile, f)
	if err != nil {
		return fmt.Errorf("Unable to write chunk to file", err)
	}

	// If this is the last chunk, merge the files
	if chunkIndex == totalChunks-1 {
		mergeChunks(totalChunks, file.Filename)
	}
	return nil
}

// 合并分片文件
func mergeChunks(totalChunks int, filename string) {
	outFilePath := filepath.Join(consts.UploadDir, filename)
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return // handle error
	}
	defer outFile.Close()

	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(consts.TempDir, "chunk_"+strconv.Itoa(i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return // handle error
		}

		io.Copy(outFile, chunkFile)
		chunkFile.Close()
		os.Remove(chunkPath) // 删除块文件
	}
}
