package service

import (
	"FileNest/internal/model"
	"FileNest/internal/service/impl"
)

// ShareService 文件分享服务接口
type ShareService interface {
	// CreateShare 创建分享
	CreateShare(filePath string, password string, expireHours int, maxDownload int, description string, creatorIP string) (*model.FileShare, error)
	// GetShareByCode 根据分享码获取分享信息
	GetShareByCode(shareCode string) (*model.FileShare, error)
	// ValidateShare 验证分享（密码、过期时间、下载次数等）
	ValidateShare(shareCode string, password string) (*model.FileShare, error)
	// DownloadSharedFile 下载分享文件
	DownloadSharedFile(shareCode string, password string) (string, error)
	// GetMyShares 获取我的分享列表
	GetMyShares(creatorIP string, page, pageSize int) ([]model.FileShare, int64, error)
	// UpdateShare 更新分享
	UpdateShare(share *model.FileShare) error
	// DeleteShare 删除分享
	DeleteShare(shareCode string) error
	// DisableShare 禁用分享
	DisableShare(shareCode string) error
	// GetShareStats 获取分享统计信息
	GetShareStats(shareCode string) (*model.FileShare, error)
}

// NewShareService 创建分享服务实例
func NewShareService() ShareService {
	return impl.NewShareServiceImpl()
}
