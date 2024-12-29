package cache

import (
	"fmt"
	"path/filepath"
)

const (
	// 缓存过期时间（秒）
	FileListExpiration  = 300  // 5分钟
	FileStatsExpiration = 300  // 5分钟
	SearchExpiration    = 60   // 1分钟
	FavoriteExpiration  = 1800 // 30分钟
)

// 文件列表缓存键
func FileListKey(path string) string {
	return fmt.Sprintf("file:list:%s", filepath.Clean(path))
}

// 文件统计信息缓存键
func FileStatsKey(path string) string {
	return fmt.Sprintf("file:stats:%s", filepath.Clean(path))
}

// 搜索结果缓存键
func SearchKey(keyword string) string {
	return fmt.Sprintf("file:search:%s", keyword)
}

// 搜索历史缓存键
const SearchHistoryKey = "file:search:history"

// 收藏夹缓存键
const FavoriteKey = "file:favorites"

// 上传进度缓存键
func UploadProgressKey(path, fileName string) string {
	fullPath := filepath.Join(path, fileName)
	return fmt.Sprintf("file:upload:progress:%s", filepath.Clean(fullPath))
}

// 获取目录相关的所有缓存键模式
func GetDirCachePatterns(path string) []string {
	path = filepath.Clean(path)
	return []string{
		fmt.Sprintf("file:list:%s*", path),
		fmt.Sprintf("file:stats:%s*", path),
	}
}

// 获取文件相关的所有缓存键
func GetFileCacheKeys(path string) []string {
	dir := filepath.Dir(path)
	return []string{
		FileListKey(dir),
		FileStatsKey(dir),
	}
}
