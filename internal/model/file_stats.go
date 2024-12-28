package model

// FileStats 文件统计信息
type FileStats struct {
	TotalFiles   int64 `json:"totalFiles"`   // 文件总数
	TotalFolders int64 `json:"totalFolders"` // 文件夹总数
	TotalSize    int64 `json:"totalSize"`    // 总大小（字节）
}
