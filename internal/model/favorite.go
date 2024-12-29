package model

import "time"

// Favorite 文件收藏
type Favorite struct {
	ID         int64     `json:"id"`         // 收藏ID
	Name       string    `json:"name"`       // 文件名
	Path       string    `json:"path"`       // 文件路径
	IsDir      bool      `json:"isDir"`      // 是否是目录
	CreateTime time.Time `json:"createTime"` // 创建时间
}
