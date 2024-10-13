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

func (h *FileServiceImpl) DownloadFile(filePath string) (string, error) {
	// 检查文件是否存在
	filePath = filepath.Join(consts.UploadDir, filePath)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("stat error: %s", err)
	}

	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist")
	}
	return absPath, nil
}

func (h *FileServiceImpl) CreateDir(path string) error {
	path = filepath.Join(consts.UploadDir, path)
	fileInfo, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat error: %s", err)
	}
	if err == nil && !fileInfo.IsDir() {
		return fmt.Errorf("path exists but is not a directory")
	}
	if err != nil {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func (h *FileServiceImpl) GetFileList(path string) ([]model.FileInfo, error) {
	p := filepath.Join(consts.UploadDir, path)
	files, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("unable to read directory: %s", err)
	}
	var fs []model.FileInfo
	for _, file := range files {
		info, e := file.Info()
		// 通过名称获取文件类型
		fileType := filepath.Ext(file.Name())
		if e != nil {
			return nil, fmt.Errorf("unable to read file info: %s", e)
		}
		fs = append(fs, model.FileInfo{
			FileName: file.Name(),
			FilePath: filepath.Join(path, file.Name()),
			FileSize: info.Size(),
			FileType: fileType,
			IsDir:    file.IsDir(),
			ModTime:  info.ModTime(),
		})
	}
	return fs, nil
}

// DeleteFile 删除文件
func (h *FileServiceImpl) DeleteFile(path string) error {
	p := filepath.Join(consts.UploadDir, path)
	fileInfo, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat error: %s", err)
	}
	if fileInfo.IsDir() {
		return os.RemoveAll(p)
	}
	return os.Remove(p)
}

func (h *FileServiceImpl) UploadFile(file *multipart.FileHeader, path, fileName string, chunkIndex, totalChunks int) error {
	err := os.MkdirAll(consts.TempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create temp directory: %s", err)
	}
	name := filepath.Base(fileName)
	tempFileName := name + "_chunk_" + strconv.Itoa(chunkIndex)
	chunkPath := filepath.Join(consts.TempDir, tempFileName)
	outFile, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("unable to create chunk file: %s", err)
	}
	defer outFile.Close()

	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open chunk file: %s", err)
	}
	defer f.Close()
	_, err = io.Copy(outFile, f)
	if err != nil {
		return fmt.Errorf("unable to write chunk to file: %s", err)
	}

	if chunkIndex == totalChunks {
		mergeChunks(totalChunks, fileName, path)
	}
	return nil
}

// 合并分片文件
func mergeChunks(totalChunks int, filename, path string) {
	outFilePath := filepath.Join(consts.UploadDir, path, filename)
	err := os.MkdirAll(filepath.Dir(outFilePath), os.ModePerm)
	if err != nil {
		glog.Errorf("unable to create upload directory %s", err)
		return
	}
	outFile, err := os.Create(outFilePath)
	if err != nil {
		glog.Errorf("unable to create output file %s", err)
		return
	}
	defer outFile.Close()

	for i := 0; i < totalChunks; i++ {
		name := filepath.Base(filename)
		tempFileName := name + "_chunk_" + strconv.Itoa(i+1)
		chunkPath := filepath.Join(consts.TempDir, tempFileName)
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			glog.Errorf("unable to open chunk file %s: %v", chunkPath, err)
			continue
		}

		_, err = io.Copy(outFile, chunkFile)
		if err != nil {
			glog.Errorf("unable to copy chunk to output file %s: %v", chunkPath, err)
			return
		}

		if err := chunkFile.Close(); err != nil {
			glog.Errorf("unable to close chunk file %s: %v", chunkPath, err)
		}
		if err := os.Remove(chunkPath); err != nil {
			glog.Errorf("unable to remove chunk file %s: %v", chunkPath, err)
		}
	}
}
