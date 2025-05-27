package model

import (
	"time"

	"gorm.io/gorm"
)

// FileShare 文件分享
type FileShare struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`                     // 分享ID
	ShareCode   string         `json:"shareCode" gorm:"type:varchar(32);uniqueIndex;not null"` // 分享码
	FilePath    string         `json:"filePath" gorm:"type:varchar(500);not null"`             // 文件路径
	FileName    string         `json:"fileName" gorm:"type:varchar(255);not null"`             // 文件名
	FileSize    int64          `json:"fileSize" gorm:"default:0"`                              // 文件大小
	IsDir       bool           `json:"isDir" gorm:"default:false"`                             // 是否是目录
	Password    string         `json:"password,omitempty" gorm:"type:varchar(32)"`             // 提取密码（可选）
	ExpireTime  *time.Time     `json:"expireTime"`                                             // 过期时间（可选）
	MaxDownload int            `json:"maxDownload" gorm:"default:0"`                           // 最大下载次数（0表示无限制）
	Downloaded  int            `json:"downloaded" gorm:"default:0"`                            // 已下载次数
	Status      string         `json:"status" gorm:"type:varchar(20);default:'active'"`        // 状态：active/expired/disabled
	Description string         `json:"description" gorm:"type:text"`                           // 分享描述
	CreatorIP   string         `json:"creatorIP" gorm:"type:varchar(45)"`                      // 创建者IP
	CreateTime  time.Time      `json:"createTime" gorm:"autoCreateTime"`                       // 创建时间
	UpdateTime  time.Time      `json:"updateTime" gorm:"autoUpdateTime"`                       // 更新时间
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                                         // 软删除
}

// TableName 指定表名
func (FileShare) TableName() string {
	return "file_shares"
}

// IsExpired 检查是否过期
func (s *FileShare) IsExpired() bool {
	if s.ExpireTime != nil && time.Now().After(*s.ExpireTime) {
		return true
	}
	return false
}

// IsMaxDownloadReached 检查是否达到最大下载次数
func (s *FileShare) IsMaxDownloadReached() bool {
	if s.MaxDownload > 0 && s.Downloaded >= s.MaxDownload {
		return true
	}
	return false
}

// CanDownload 检查是否可以下载
func (s *FileShare) CanDownload() bool {
	return s.Status == "active" && !s.IsExpired() && !s.IsMaxDownloadReached()
}

// IncrementDownload 增加下载次数
func (s *FileShare) IncrementDownload() {
	s.Downloaded++
	if s.MaxDownload > 0 && s.Downloaded >= s.MaxDownload {
		s.Status = "expired"
	}
}
