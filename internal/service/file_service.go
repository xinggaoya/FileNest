package service

import (
	"FileNest/internal/model"
	"FileNest/internal/service/impl"
)

/**
  @author: XingGao
  @date: 2024/9/28
**/

type FileService interface {
	GetFileList(path string) ([]model.FileInfo, error)
	UploadFile(path, fileName string, totalChunks int) error
	CreateDir(path string) error
	DeleteFile(path string) error
	// DownloadFile 下载
	DownloadFile(path string) (string, error)
}

func NewFileService() FileService {
	return &impl.FileServiceImpl{}
}
