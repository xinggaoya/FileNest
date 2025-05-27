package repository

import (
	"FileNest/common/database"
	"FileNest/internal/model"
	"fmt"

	"gorm.io/gorm"
)

// ShareRepository 分享仓储接口
type ShareRepository interface {
	Create(share *model.FileShare) error
	GetByCode(shareCode string) (*model.FileShare, error)
	Update(share *model.FileShare) error
	DeleteByCode(shareCode string) error
	GetByCreatorIP(creatorIP string, page, pageSize int) ([]model.FileShare, int64, error)
	GetExpiredShares() ([]model.FileShare, error)
	CleanupExpiredShares() error
	ExistsByCode(shareCode string) (bool, error)
}

// shareRepository 分享仓储实现
type shareRepository struct {
	db *gorm.DB
}

// NewShareRepository 创建分享仓储实例
func NewShareRepository() ShareRepository {
	return &shareRepository{
		db: database.GetDB(),
	}
}

// Create 创建分享
func (r *shareRepository) Create(share *model.FileShare) error {
	// 检查分享码是否已存在
	var existing model.FileShare
	err := r.db.Where("share_code = ?", share.ShareCode).First(&existing).Error
	if err == nil {
		return fmt.Errorf("分享码已存在")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("检查分享码失败: %v", err)
	}

	// 创建新分享
	if err := r.db.Create(share).Error; err != nil {
		return fmt.Errorf("创建分享失败: %v", err)
	}
	return nil
}

// GetByCode 根据分享码获取分享
func (r *shareRepository) GetByCode(shareCode string) (*model.FileShare, error) {
	var share model.FileShare
	err := r.db.Where("share_code = ?", shareCode).First(&share).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("分享不存在")
		}
		return nil, fmt.Errorf("查询分享失败: %v", err)
	}
	return &share, nil
}

// Update 更新分享
func (r *shareRepository) Update(share *model.FileShare) error {
	err := r.db.Save(share).Error
	if err != nil {
		return fmt.Errorf("更新分享失败: %v", err)
	}
	return nil
}

// DeleteByCode 根据分享码删除分享
func (r *shareRepository) DeleteByCode(shareCode string) error {
	result := r.db.Where("share_code = ?", shareCode).Delete(&model.FileShare{})
	if result.Error != nil {
		return fmt.Errorf("删除分享失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("分享不存在")
	}
	return nil
}

// GetByCreatorIP 根据创建者IP分页获取分享
func (r *shareRepository) GetByCreatorIP(creatorIP string, page, pageSize int) ([]model.FileShare, int64, error) {
	var shares []model.FileShare
	var total int64

	// 获取总数
	if err := r.db.Model(&model.FileShare{}).Where("creator_ip = ?", creatorIP).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询分享总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Where("creator_ip = ?", creatorIP).
		Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&shares).Error

	if err != nil {
		return nil, 0, fmt.Errorf("分页查询分享失败: %v", err)
	}

	return shares, total, nil
}

// GetExpiredShares 获取已过期的分享
func (r *shareRepository) GetExpiredShares() ([]model.FileShare, error) {
	var shares []model.FileShare
	err := r.db.Where("expire_time IS NOT NULL AND expire_time < NOW()").
		Or("status = ?", "expired").
		Find(&shares).Error

	if err != nil {
		return nil, fmt.Errorf("查询过期分享失败: %v", err)
	}
	return shares, nil
}

// CleanupExpiredShares 清理过期的分享
func (r *shareRepository) CleanupExpiredShares() error {
	result := r.db.Where("expire_time IS NOT NULL AND expire_time < NOW()").
		Or("status = ?", "expired").
		Delete(&model.FileShare{})

	if result.Error != nil {
		return fmt.Errorf("清理过期分享失败: %v", result.Error)
	}
	return nil
}

// ExistsByCode 检查分享码是否存在
func (r *shareRepository) ExistsByCode(shareCode string) (bool, error) {
	var count int64
	err := r.db.Model(&model.FileShare{}).Where("share_code = ?", shareCode).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查分享码失败: %v", err)
	}
	return count > 0, nil
}
