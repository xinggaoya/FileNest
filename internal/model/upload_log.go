package model

import (
	"time"

	"gorm.io/gorm"
)

// UploadLog 上传日志
type UploadLog struct {
	ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`                  // 上传ID
	FileName       string         `json:"fileName" gorm:"type:varchar(255);not null"`          // 文件名
	FilePath       string         `json:"filePath" gorm:"type:varchar(500);not null"`          // 文件路径
	FileSize       int64          `json:"fileSize" gorm:"default:0"`                           // 文件大小
	FileType       string         `json:"fileType" gorm:"type:varchar(50)"`                    // 文件类型
	UploadType     string         `json:"uploadType" gorm:"type:varchar(20);default:'normal'"` // 上传类型：normal/chunk
	Status         string         `json:"status" gorm:"type:varchar(20);default:'pending'"`    // 状态：pending/uploading/success/failed
	TotalChunks    int            `json:"totalChunks" gorm:"default:0"`                        // 总分块数
	UploadedChunks int            `json:"uploadedChunks" gorm:"default:0"`                     // 已上传分块数
	Progress       float64        `json:"progress" gorm:"type:decimal(5,2);default:0"`         // 上传进度
	ErrorMessage   string         `json:"errorMessage" gorm:"type:text"`                       // 错误信息
	ClientIP       string         `json:"clientIP" gorm:"type:varchar(45)"`                    // 客户端IP
	UserAgent      string         `json:"userAgent" gorm:"type:text"`                          // 用户代理
	CreateTime     time.Time      `json:"createTime" gorm:"autoCreateTime"`                    // 创建时间
	UpdateTime     time.Time      `json:"updateTime" gorm:"autoUpdateTime"`                    // 更新时间
	CompletedAt    *time.Time     `json:"completedAt"`                                         // 完成时间
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`                                      // 软删除
}

// TableName 指定表名
func (UploadLog) TableName() string {
	return "file_upload_logs"
}

// IsCompleted 检查是否完成上传
func (u *UploadLog) IsCompleted() bool {
	return u.Status == "success"
}

// IsFailed 检查是否上传失败
func (u *UploadLog) IsFailed() bool {
	return u.Status == "failed"
}

// UpdateProgress 更新上传进度
func (u *UploadLog) UpdateProgress() {
	if u.TotalChunks > 0 {
		u.Progress = float64(u.UploadedChunks) / float64(u.TotalChunks) * 100
	}
}
