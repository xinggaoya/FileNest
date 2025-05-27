package repository

import (
	"FileNest/common/database"
	"FileNest/internal/model"
	"fmt"

	"gorm.io/gorm"
)

// FavoriteRepository 收藏仓储接口
type FavoriteRepository interface {
	Create(favorite *model.Favorite) error
	GetByPath(path string) (*model.Favorite, error)
	GetAll() ([]model.Favorite, error)
	GetPaginated(page, pageSize int) ([]model.Favorite, int64, error)
	Update(favorite *model.Favorite) error
	DeleteByPath(path string) error
	DeleteByID(id int64) error
	ExistsByPath(path string) (bool, error)
	GetRecentlyAdded(limit int) ([]model.Favorite, error)
	Search(keyword string) ([]model.Favorite, error)
}

// favoriteRepository 收藏仓储实现
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓储实例
func NewFavoriteRepository() FavoriteRepository {
	return &favoriteRepository{
		db: database.GetDB(),
	}
}

// Create 创建收藏
func (r *favoriteRepository) Create(favorite *model.Favorite) error {
	// 检查是否已存在
	var existing model.Favorite
	err := r.db.Where("path = ?", favorite.Path).First(&existing).Error
	if err == nil {
		return fmt.Errorf("文件已在收藏列表中")
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("检查收藏状态失败: %v", err)
	}

	// 创建新收藏
	if err := r.db.Create(favorite).Error; err != nil {
		return fmt.Errorf("添加收藏失败: %v", err)
	}
	return nil
}

// GetByPath 根据路径获取收藏
func (r *favoriteRepository) GetByPath(path string) (*model.Favorite, error) {
	var favorite model.Favorite
	err := r.db.Where("path = ?", path).First(&favorite).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("收藏不存在")
		}
		return nil, fmt.Errorf("查询收藏失败: %v", err)
	}
	return &favorite, nil
}

// GetAll 获取所有收藏
func (r *favoriteRepository) GetAll() ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := r.db.Order("create_time DESC").Find(&favorites).Error
	if err != nil {
		return nil, fmt.Errorf("查询收藏列表失败: %v", err)
	}
	return favorites, nil
}

// GetPaginated 分页获取收藏
func (r *favoriteRepository) GetPaginated(page, pageSize int) ([]model.Favorite, int64, error) {
	var favorites []model.Favorite
	var total int64

	// 获取总数
	if err := r.db.Model(&model.Favorite{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询收藏总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Order("create_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&favorites).Error

	if err != nil {
		return nil, 0, fmt.Errorf("分页查询收藏失败: %v", err)
	}

	return favorites, total, nil
}

// Update 更新收藏
func (r *favoriteRepository) Update(favorite *model.Favorite) error {
	err := r.db.Save(favorite).Error
	if err != nil {
		return fmt.Errorf("更新收藏失败: %v", err)
	}
	return nil
}

// DeleteByPath 根据路径删除收藏
func (r *favoriteRepository) DeleteByPath(path string) error {
	result := r.db.Where("path = ?", path).Delete(&model.Favorite{})
	if result.Error != nil {
		return fmt.Errorf("删除收藏失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("收藏不存在")
	}
	return nil
}

// DeleteByID 根据ID删除收藏
func (r *favoriteRepository) DeleteByID(id int64) error {
	result := r.db.Delete(&model.Favorite{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除收藏失败: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("收藏不存在")
	}
	return nil
}

// ExistsByPath 检查路径是否已收藏
func (r *favoriteRepository) ExistsByPath(path string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("path = ?", path).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查收藏状态失败: %v", err)
	}
	return count > 0, nil
}

// GetRecentlyAdded 获取最近添加的收藏
func (r *favoriteRepository) GetRecentlyAdded(limit int) ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := r.db.Order("create_time DESC").Limit(limit).Find(&favorites).Error
	if err != nil {
		return nil, fmt.Errorf("查询最近收藏失败: %v", err)
	}
	return favorites, nil
}

// Search 搜索收藏
func (r *favoriteRepository) Search(keyword string) ([]model.Favorite, error) {
	var favorites []model.Favorite
	searchPattern := "%" + keyword + "%"
	err := r.db.Where("name LIKE ? OR path LIKE ?", searchPattern, searchPattern).
		Order("create_time DESC").
		Find(&favorites).Error

	if err != nil {
		return nil, fmt.Errorf("搜索收藏失败: %v", err)
	}
	return favorites, nil
}
