package model

import "time"

// Favorite 文件收藏
type Favorite struct {
	ID        int64     `json:"id"`        // 收藏ID
	FilePath  string    `json:"filePath"`  // 文件路径
	FileName  string    `json:"fileName"`  // 文件名
	IsDir     bool      `json:"isDir"`     // 是否是目录
	CreatedAt time.Time `json:"createdAt"` // 创建时间
}
