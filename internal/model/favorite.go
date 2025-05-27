package model

import (
	"time"

	"gorm.io/gorm"
)

// Favorite 文件收藏
type Favorite struct {
	ID         int64          `json:"id" gorm:"primaryKey;autoIncrement"`                 // 收藏ID
	Name       string         `json:"name" gorm:"type:varchar(255);not null"`             // 文件名
	Path       string         `json:"path" gorm:"type:varchar(500);not null;uniqueIndex"` // 文件路径
	IsDir      bool           `json:"isDir" gorm:"default:false"`                         // 是否是目录
	FileSize   int64          `json:"fileSize" gorm:"default:0"`                          // 文件大小
	FileType   string         `json:"fileType" gorm:"type:varchar(50)"`                   // 文件类型
	CreateTime time.Time      `json:"createTime" gorm:"autoCreateTime"`                   // 创建时间
	UpdateTime time.Time      `json:"updateTime" gorm:"autoUpdateTime"`                   // 更新时间
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`                                     // 软删除
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "file_favorites"
}
