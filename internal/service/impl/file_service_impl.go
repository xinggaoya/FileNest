package impl

import (
	"FileNest/common/glog"
	"FileNest/internal/consts"
	"FileNest/internal/model"
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

type FileServiceImpl struct{}

func (h *FileServiceImpl) GetFileList(path string) ([]model.FileInfo, error) {
	p := filepath.Join(consts.UploadDir, path)
	// 读取路径下的文件列表
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory", err)
	}
	var fs []model.FileInfo

	for _, file := range files {
		info, e := file.Info()
		if e != nil {
			return nil, fmt.Errorf("unable to read file info", e)
		}
		fs = append(fs, model.FileInfo{
			FileName: file.Name(),
			FilePath: filepath.Join(path, file.Name()),
			FileSize: info.Size(),
			FileType: file.Type().String(),
			IsDir:    file.IsDir(),
			ModTime:  info.ModTime(),
		})
	}

	return fs, err
}

func (h *FileServiceImpl) UploadFile(file *multipart.FileHeader, chunkIndex, totalChunks int) error {
	// 判断文件夹是否存在，不存在则创建
	err := os.MkdirAll(consts.TempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create temp directory", err)
	}
	// Create the file to store the chunk
	chunkPath := filepath.Join(consts.TempDir, "chunk_"+strconv.Itoa(chunkIndex))
	outFile, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("unable to create chunk file", err)
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
	if chunkIndex == totalChunks {
		mergeChunks(totalChunks, file.Filename)
	}
	return nil
}

// 合并分片文件
func mergeChunks(totalChunks int, filename string) {
	// 判断文件夹是否存在，不存在则创建
	err := os.MkdirAll(consts.UploadDir, os.ModePerm)
	if err != nil {
		glog.Errorf("unable to create upload directory %s", err)
		return
	}
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
