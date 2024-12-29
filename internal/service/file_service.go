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
	UploadFile(path, fileName string, totalChunks int, override bool) error
	CreateDir(path string) error
	DeleteFile(path string, force bool) error
	// DownloadFile 下载
	DownloadFile(path string) (string, error)
	CreateFolder(path string) error
	RemoveFile(path string, force bool) error
	GetFileStats(path string) (*model.FileStats, error)
	// SearchFiles 搜索文件
	SearchFiles(keyword string) ([]model.FileInfo, error)
	// AddFavorite 添加收藏
	AddFavorite(filePath string) error
	// RemoveFavorite 取消收藏
	RemoveFavorite(filePath string) error
	// GetFavorites 获取收藏列表
	GetFavorites() ([]model.Favorite, error)
}

func NewFileService() FileService {
	return &impl.FileServiceImpl{}
}
