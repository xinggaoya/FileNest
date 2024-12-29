package impl

import (
	"FileNest/common/glog"
	"FileNest/internal/cache"
	"FileNest/internal/consts"
	"FileNest/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

/**
  @author: XingGao
  @date: 2024/9/28
**/

type FileServiceImpl struct{}

func init() {
	// 确保上传目录存在
	if err := os.MkdirAll(consts.UploadDir, os.ModePerm); err != nil {
		glog.Errorf("创建上传目录失败: %s", err)
		panic(err)
	}
}

func (h *FileServiceImpl) DownloadFile(filePath string) (string, error) {
	// 检查文件是否存在
	filePath = filepath.Join(consts.UploadDir, filePath)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("stat error: %s", err)
	}

	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist")
	}
	return absPath, nil
}

func (h *FileServiceImpl) CreateDir(path string) error {
	path = filepath.Join(consts.UploadDir, path)
	fileInfo, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("stat error: %s", err)
	}
	if err == nil && !fileInfo.IsDir() {
		return fmt.Errorf("path exists but is not a directory")
	}
	if err != nil {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// GetFileList 获取文件列表（带缓存）
func (s *FileServiceImpl) GetFileList(path string) ([]model.FileInfo, error) {
	glog.Infof("开始获取文件列表，路径: %s", path)

	// 尝试从缓存获取
	cacheKey := cache.FileListKey(path)
	if cached, err := cache.Get(cacheKey); err == nil {
		var files []model.FileInfo
		if err := json.Unmarshal([]byte(cached), &files); err == nil {
			glog.Infof("从缓存获取文件列表成功，路径: %s", path)
			return files, nil
		}
	}

	// 从文件系统获取
	files, err := s.getFileListFromFS(path)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	if cacheData, err := json.Marshal(files); err == nil {
		cache.Set(cacheKey, cacheData, time.Duration(cache.FileListExpiration)*time.Second)
	}

	return files, nil
}

// getFileListFromFS 从文件系统获取文件列表
func (s *FileServiceImpl) getFileListFromFS(path string) ([]model.FileInfo, error) {
	absPath := filepath.Join(consts.UploadDir, path)
	glog.Infof("从文件系统获取文件列表，完整路径: %s", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("路径不存在: %s", path)
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %s", err)
	}

	var files []model.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			glog.Errorf("获取文件信息失败: %s", err)
			continue
		}

		fileInfo := model.FileInfo{
			FileName: entry.Name(),
			FilePath: filepath.Join(path, entry.Name()),
			FileSize: info.Size(),
			FileType: filepath.Ext(entry.Name()),
			IsDir:    entry.IsDir(),
			ModTime:  info.ModTime().Format(time.DateTime),
		}
		files = append(files, fileInfo)
	}

	return files, nil
}

// GetFileStats 获取文件统计信息（带缓存）
func (s *FileServiceImpl) GetFileStats(path string) (*model.FileStats, error) {
	glog.Infof("开始获取文件统计信息，路径: %s", path)

	// 尝试从缓存获取
	cacheKey := cache.FileStatsKey(path)
	if cached, err := cache.Get(cacheKey); err == nil {
		var stats model.FileStats
		if err := json.Unmarshal([]byte(cached), &stats); err == nil {
			glog.Infof("从缓存获取文件统计信息成功，路径: %s", path)
			return &stats, nil
		}
	}

	// 从文件系统获取
	stats, err := s.getFileStatsFromFS(path)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	if cacheData, err := json.Marshal(stats); err == nil {
		cache.Set(cacheKey, cacheData, time.Duration(cache.FileStatsExpiration)*time.Second)
	}

	return stats, nil
}

// getFileStatsFromFS 从文件系统获取文件统计信息
func (s *FileServiceImpl) getFileStatsFromFS(path string) (*model.FileStats, error) {
	absPath := filepath.Join(consts.UploadDir, path)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("路径不存在: %s", path)
	}

	stats := &model.FileStats{}
	err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path != absPath {
				stats.TotalFolders++
			}
		} else {
			stats.TotalFiles++
			stats.TotalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("统计文件信息失败: %s", err)
	}

	return stats, nil
}

// SearchFiles 搜索文件（带缓存）
func (s *FileServiceImpl) SearchFiles(keyword string) ([]model.FileInfo, error) {
	glog.Infof("开始搜索文件，关键词: %s", keyword)

	ctx := context.Background()
	redisClient := cache.GetRedisClient()

	// 记录搜索历史
	score := float64(time.Now().Unix())
	redisClient.ZAdd(ctx, cache.SearchHistoryKey, redis.Z{
		Score:  score,
		Member: keyword,
	})
	// 只保留最近 10 条记录
	redisClient.ZRemRangeByRank(ctx, cache.SearchHistoryKey, 0, -11)

	// 尝试从缓存获取搜索结果
	cacheKey := cache.SearchKey(keyword)
	if cached, err := cache.Get(cacheKey); err == nil {
		var files []model.FileInfo
		if err := json.Unmarshal([]byte(cached), &files); err == nil {
			glog.Infof("从缓存获取搜索结果成功，关键词: %s", keyword)
			return files, nil
		}
	}

	// 执行搜索
	files, err := s.searchFilesInFS(keyword)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	if cacheData, err := json.Marshal(files); err == nil {
		cache.Set(cacheKey, cacheData, time.Duration(cache.SearchExpiration)*time.Second)
	}

	return files, nil
}

// searchFilesInFS 在文件系统中搜索文件
func (s *FileServiceImpl) searchFilesInFS(keyword string) ([]model.FileInfo, error) {
	var files []model.FileInfo
	err := filepath.Walk(consts.UploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查文件名是否匹配关键词
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(keyword)) {
			relPath, err := filepath.Rel(consts.UploadDir, path)
			if err != nil {
				return err
			}

			fileInfo := model.FileInfo{
				FileName: info.Name(),
				FilePath: relPath,
				FileSize: info.Size(),
				FileType: filepath.Ext(info.Name()),
				IsDir:    info.IsDir(),
				ModTime:  info.ModTime().Format(time.DateTime),
			}
			files = append(files, fileInfo)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("搜索文件失败: %s", err)
	}

	return files, nil
}

// DeleteFile 删除文件（清除相关缓存）
func (s *FileServiceImpl) DeleteFile(path string, force bool) error {
	glog.Infof("开始删除文件，路径: %s, 强制删除: %v", path, force)

	// 删除文件
	err := s.deleteFileFromFS(path, force)
	if err != nil {
		return err
	}

	// 清除相关缓存
	s.clearFileRelatedCache(path)
	return nil
}

// deleteFileFromFS 从文件系统删除文件
func (s *FileServiceImpl) deleteFileFromFS(path string, force bool) error {
	path = filepath.Clean(path)
	fullPath := filepath.Join(consts.UploadDir, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("检查路径状态失败: %s", err)
	}

	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return fmt.Errorf("读取目录失败: %s", err)
		}
		if len(entries) > 0 && !force {
			return fmt.Errorf("目录不为空，如需删除请勾选\"强制删除\"选项")
		}
		if force {
			err = os.RemoveAll(fullPath)
		} else {
			err = os.Remove(fullPath)
		}
	} else {
		err = os.Remove(fullPath)
	}

	if err != nil {
		return fmt.Errorf("删除失败: %s", err)
	}

	return nil
}

// clearFileRelatedCache 清除文件相关的所有缓存
func (s *FileServiceImpl) clearFileRelatedCache(path string) {
	// 获取需要清除的缓存键
	cacheKeys := cache.GetFileCacheKeys(path)

	// 清除直接相关的缓存
	if len(cacheKeys) > 0 {
		cache.Del(cacheKeys...)
	}

	// 清除目录相关的缓存模式
	patterns := cache.GetDirCachePatterns(path)
	for _, pattern := range patterns {
		cache.DelByPattern(pattern)
	}
}

// UploadFile 上传文件（更新缓存）
func (s *FileServiceImpl) UploadFile(path, fileName string, totalChunks int, override bool) error {
	glog.Infof("开始上传文件，路径: %s, 文件名: %s", path, fileName)

	// 设置上传进度
	progressKey := cache.UploadProgressKey(path, fileName)
	cache.HSet(progressKey,
		"status", "uploading",
		"progress", "0",
		"total_chunks", totalChunks,
	)
	cache.Expire(progressKey, time.Hour)

	// 执行上传
	err := s.uploadFileToFS(path, fileName, override)
	if err != nil {
		// 更新失败状态
		cache.HSet(progressKey, "status", "error", "error", err.Error())
		return err
	}

	// 清除相关缓存
	s.clearFileRelatedCache(filepath.Join(path, fileName))
	return nil
}

// uploadFileToFS 上传文件到文件系统
func (s *FileServiceImpl) uploadFileToFS(path, fileName string, override bool) error {
	path = filepath.Clean(path)
	if path == "." {
		path = ""
	}

	outFilePath := filepath.Join(consts.UploadDir, path, fileName)
	targetDir := filepath.Dir(outFilePath)

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目标目录失败: %s", err)
	}

	if err := s.checkDirectoryWritePermission(targetDir); err != nil {
		return fmt.Errorf("目标目录无写入权限: %s", err)
	}

	if _, err := os.Stat(outFilePath); err == nil {
		if !override {
			return fmt.Errorf("文件已存在: %s", fileName)
		}
		if err := os.Remove(outFilePath); err != nil {
			return fmt.Errorf("删除已存在文件失败: %s", err)
		}
	}

	return nil
}

// checkDirectoryWritePermission 检查目录是否有写入权限
func (s *FileServiceImpl) checkDirectoryWritePermission(dir string) error {
	tempFile := filepath.Join(dir, ".write_test")
	f, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	f.Close()
	return os.Remove(tempFile)
}

// AddFavorite 添加收藏
func (s *FileServiceImpl) AddFavorite(filePath string) error {
	glog.Infof("添加收藏，文件路径: %s", filePath)

	// 检查文件是否存在
	absPath := filepath.Join(consts.UploadDir, filePath)
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("文件不存在: %s", filePath)
		}
		return fmt.Errorf("获取文件信息失败: %s", err)
	}

	// 创建收藏记录
	favorite := model.Favorite{
		Name:       info.Name(),
		Path:       filePath,
		IsDir:      info.IsDir(),
		CreateTime: time.Now(),
	}

	// TODO: 保存到数据库
	glog.Infof("收藏成功: %+v", favorite)
	return nil
}

// RemoveFavorite 取消收藏
func (s *FileServiceImpl) RemoveFavorite(filePath string) error {
	glog.Infof("取消收藏，文件路径: %s", filePath)
	// TODO: 从数据库中删除
	return nil
}

// GetFavorites 获取收藏列表
func (s *FileServiceImpl) GetFavorites() ([]model.Favorite, error) {
	glog.Info("获取收藏列表")
	// TODO: 从数据库中查询
	return []model.Favorite{}, nil
}

// CreateFolder 创建文件夹
func (h *FileServiceImpl) CreateFolder(path string) error {
	glog.Infof("开始创建文件夹，路径: %s", path)

	// 规范化路径
	path = filepath.Clean(path)
	if path == "." {
		path = ""
	}

	// 构建完整的文件夹路径
	folderPath := filepath.Join(consts.UploadDir, path)
	glog.Infof("目标文件夹路径: %s", folderPath)

	// 检查路径是否已存在
	if info, err := os.Stat(folderPath); err == nil {
		if info.IsDir() {
			glog.Errorf("文件夹已存在: %s", folderPath)
			return fmt.Errorf("文件夹已存在: %s", path)
		}
		glog.Errorf("路径已存在但不是文件夹: %s", folderPath)
		return fmt.Errorf("路径已存在但不是文件夹: %s", path)
	} else if !os.IsNotExist(err) {
		glog.Errorf("检查路径状态失败: %s", err)
		return fmt.Errorf("检查路径状态失败: %s", err)
	}

	// 创建文件夹
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		glog.Errorf("创建文件夹失败: %s, 路径: %s", err, folderPath)
		return fmt.Errorf("创建文件夹失败: %s", err)
	}

	glog.Infof("文件夹创建成功: %s", folderPath)
	return nil
}

// RemoveFile 删除文件（实际调用 DeleteFile 方法）
func (h *FileServiceImpl) RemoveFile(path string, force bool) error {
	return h.DeleteFile(path, force)
}

// RenameFile 重命名文件或文件夹
func (h *FileServiceImpl) RenameFile(oldPath string, newName string) error {
	glog.Infof("开始重命名，原路径: %s, 新名称: %s", oldPath, newName)

	// 构建完整的原路径
	oldFullPath := filepath.Join(consts.UploadDir, oldPath)

	// 获取父目录
	parentDir := filepath.Dir(oldPath)
	// 构建新路径
	newPath := filepath.Join(parentDir, newName)
	newFullPath := filepath.Join(consts.UploadDir, newPath)

	// 检查新路径是否已存在
	if _, err := os.Stat(newFullPath); err == nil {
		return fmt.Errorf("目标路径已存在: %s", newPath)
	}

	// 执行重命名
	if err := os.Rename(oldFullPath, newFullPath); err != nil {
		glog.Errorf("重命名失败: %s", err)
		return fmt.Errorf("重命名失败: %s", err)
	}

	glog.Infof("重命名成功: %s -> %s", oldPath, newPath)
	return nil
}

// CopyFile 复制文件或文件夹
func (h *FileServiceImpl) CopyFile(srcPath string, destPath string) error {
	glog.Infof("开始复制，源路径: %s, 目标路径: %s", srcPath, destPath)

	srcFullPath := filepath.Join(consts.UploadDir, srcPath)
	destFullPath := filepath.Join(consts.UploadDir, destPath)

	// 检查源路径是否存在
	srcInfo, err := os.Stat(srcFullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("源文件不存在: %s", srcPath)
		}
		return fmt.Errorf("获取源文件信息失败: %s", err)
	}

	// 检查目标路径是否已存在
	if _, err := os.Stat(destFullPath); err == nil {
		return fmt.Errorf("目标路径已存在: %s", destPath)
	}

	if srcInfo.IsDir() {
		// 复制目录
		if err := h.copyDir(srcFullPath, destFullPath); err != nil {
			return fmt.Errorf("复制目录失败: %s", err)
		}
	} else {
		// 复制文件
		if err := h.copyFileContent(srcFullPath, destFullPath); err != nil {
			return fmt.Errorf("复制文件失败: %s", err)
		}
	}

	glog.Infof("复制成功: %s -> %s", srcPath, destPath)
	return nil
}

// MoveFile 移动文件或文件夹
func (h *FileServiceImpl) MoveFile(srcPath string, destPath string) error {
	glog.Infof("开始移动，源路径: %s, 目标路径: %s", srcPath, destPath)

	srcFullPath := filepath.Join(consts.UploadDir, srcPath)
	destFullPath := filepath.Join(consts.UploadDir, destPath)

	// 检查源路径是否存在
	if _, err := os.Stat(srcFullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("源文件不存在: %s", srcPath)
		}
		return fmt.Errorf("获取源文件信息失败: %s", err)
	}

	// 检查目标路径是否已存在
	if _, err := os.Stat(destFullPath); err == nil {
		return fmt.Errorf("目标路径已存在: %s", destPath)
	}

	// 确保目标目录存在
	destDir := filepath.Dir(destFullPath)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目标目录失败: %s", err)
	}

	// 执行移动
	if err := os.Rename(srcFullPath, destFullPath); err != nil {
		glog.Errorf("移动失败: %s", err)
		return fmt.Errorf("移动失败: %s", err)
	}

	glog.Infof("移动成功: %s -> %s", srcPath, destPath)
	return nil
}

// copyDir 复制目录
func (h *FileServiceImpl) copyDir(src string, dest string) error {
	// 创建目标目录
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	// 读取源目录
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			if err := h.copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// 复制文件
			if err := h.copyFileContent(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFileContent 复制文件内容
func (h *FileServiceImpl) copyFileContent(src string, dest string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制内容
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	// 保持文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dest, srcInfo.Mode())
}
