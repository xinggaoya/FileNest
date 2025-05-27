package impl

import (
	"FileNest/common/glog"
	"FileNest/internal/consts"
	"FileNest/internal/model"
	"FileNest/internal/repository"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ShareServiceImpl 文件分享服务实现
type ShareServiceImpl struct {
	shareRepo repository.ShareRepository
}

// NewShareServiceImpl 创建分享服务实例
func NewShareServiceImpl() *ShareServiceImpl {
	return &ShareServiceImpl{
		shareRepo: repository.NewShareRepository(),
	}
}

// CreateShare 创建分享
func (s *ShareServiceImpl) CreateShare(filePath string, password string, expireHours int, maxDownload int, description string, creatorIP string) (*model.FileShare, error) {
	glog.Infof("创建分享，文件路径: %s, 创建者IP: %s", filePath, creatorIP)

	// 检查文件是否存在
	absPath := filepath.Join(consts.UploadDir, filePath)
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件不存在: %s", filePath)
		}
		return nil, fmt.Errorf("获取文件信息失败: %s", err)
	}

	// 生成唯一的分享码
	shareCode, err := s.generateShareCode()
	if err != nil {
		return nil, fmt.Errorf("生成分享码失败: %v", err)
	}

	// 计算过期时间
	var expireTime *time.Time
	if expireHours > 0 {
		expire := time.Now().Add(time.Duration(expireHours) * time.Hour)
		expireTime = &expire
	}

	// 创建分享记录
	share := &model.FileShare{
		ShareCode:   shareCode,
		FilePath:    filePath,
		FileName:    info.Name(),
		FileSize:    info.Size(),
		IsDir:       info.IsDir(),
		Password:    password,
		ExpireTime:  expireTime,
		MaxDownload: maxDownload,
		Downloaded:  0,
		Status:      "active",
		Description: description,
		CreatorIP:   creatorIP,
	}

	// 保存到数据库
	if err := s.shareRepo.Create(share); err != nil {
		return nil, fmt.Errorf("创建分享失败: %v", err)
	}

	glog.Infof("分享创建成功，分享码: %s", shareCode)
	return share, nil
}

// GetShareByCode 根据分享码获取分享信息
func (s *ShareServiceImpl) GetShareByCode(shareCode string) (*model.FileShare, error) {
	share, err := s.shareRepo.GetByCode(shareCode)
	if err != nil {
		return nil, err
	}

	// 检查分享状态
	if share.IsExpired() {
		share.Status = "expired"
		s.shareRepo.Update(share)
		return nil, fmt.Errorf("分享已过期")
	}

	return share, nil
}

// ValidateShare 验证分享（密码、过期时间、下载次数等）
func (s *ShareServiceImpl) ValidateShare(shareCode string, password string) (*model.FileShare, error) {
	share, err := s.GetShareByCode(shareCode)
	if err != nil {
		return nil, err
	}

	// 检查密码
	if share.Password != "" && share.Password != password {
		return nil, fmt.Errorf("提取密码错误")
	}

	// 检查是否可以下载
	if !share.CanDownload() {
		if share.IsExpired() {
			return nil, fmt.Errorf("分享已过期")
		}
		if share.IsMaxDownloadReached() {
			return nil, fmt.Errorf("下载次数已达上限")
		}
		if share.Status != "active" {
			return nil, fmt.Errorf("分享已被禁用")
		}
	}

	return share, nil
}

// DownloadSharedFile 下载分享文件
func (s *ShareServiceImpl) DownloadSharedFile(shareCode string, password string) (string, error) {
	// 验证分享
	share, err := s.ValidateShare(shareCode, password)
	if err != nil {
		return "", err
	}

	// 检查文件是否存在
	absPath := filepath.Join(consts.UploadDir, share.FilePath)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("分享的文件不存在")
	}

	// 增加下载次数
	share.IncrementDownload()
	if err := s.shareRepo.Update(share); err != nil {
		glog.Errorf("更新下载次数失败: %v", err)
	}

	glog.Infof("分享文件下载，分享码: %s, 文件: %s", shareCode, share.FileName)
	return absPath, nil
}

// GetMyShares 获取我的分享列表
func (s *ShareServiceImpl) GetMyShares(creatorIP string, page, pageSize int) ([]model.FileShare, int64, error) {
	return s.shareRepo.GetByCreatorIP(creatorIP, page, pageSize)
}

// UpdateShare 更新分享
func (s *ShareServiceImpl) UpdateShare(share *model.FileShare) error {
	return s.shareRepo.Update(share)
}

// DeleteShare 删除分享
func (s *ShareServiceImpl) DeleteShare(shareCode string) error {
	return s.shareRepo.DeleteByCode(shareCode)
}

// DisableShare 禁用分享
func (s *ShareServiceImpl) DisableShare(shareCode string) error {
	share, err := s.shareRepo.GetByCode(shareCode)
	if err != nil {
		return err
	}

	share.Status = "disabled"
	return s.shareRepo.Update(share)
}

// GetShareStats 获取分享统计信息
func (s *ShareServiceImpl) GetShareStats(shareCode string) (*model.FileShare, error) {
	return s.shareRepo.GetByCode(shareCode)
}

// generateShareCode 生成唯一的分享码
func (s *ShareServiceImpl) generateShareCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8

	for attempts := 0; attempts < 10; attempts++ {
		// 生成随机字符串
		b := make([]byte, codeLength)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}

		var result strings.Builder
		for _, byteVal := range b {
			result.WriteByte(charset[int(byteVal)%len(charset)])
		}

		shareCode := result.String()

		// 检查是否已存在
		exists, err := s.shareRepo.ExistsByCode(shareCode)
		if err != nil {
			return "", err
		}

		if !exists {
			return shareCode, nil
		}
	}

	return "", fmt.Errorf("生成唯一分享码失败")
}
