package service

import (
	"FileNest/internal/model"
	"FileNest/internal/service/impl"
	"mime/multipart"
)

/**
  @author: XingGao
  @date: 2024/9/28
**/

type FileService interface {
	GetFileList(path string) ([]model.FileInfo, error)
	UploadFile(file *multipart.FileHeader, chunkIndex, totalChunks int) error
}

func NewFileService() FileService {
	return &impl.FileServiceImpl{}
}
