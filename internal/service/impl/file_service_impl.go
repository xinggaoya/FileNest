package impl

import (
	"FileNest/common/glog"
	"FileNest/internal/consts"
	"FileNest/internal/model"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// GetFileList 获取文件列表
func (s *FileServiceImpl) GetFileList(path string) ([]model.FileInfo, error) {
	glog.Infof("开始获取文件列表，路径: %s", path)

	// 确保路径存在
	absPath := filepath.Join(consts.UploadDir, path)
	glog.Infof("完整路径: %s", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		glog.Errorf("路径不存在: %s", absPath)
		return nil, fmt.Errorf("路径不存在: %s", path)
	}

	// 读取目录内容
	entries, err := os.ReadDir(absPath)
	if err != nil {
		glog.Errorf("读取目录失败: %s", err)
		return nil, fmt.Errorf("读取目录失败: %s", err)
	}

	// 构建文件列表
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

	glog.Infof("成功获取文件列表，共 %d 个文件", len(files))
	return files, nil
}

// DeleteFile 删除文件或文件夹
func (h *FileServiceImpl) DeleteFile(path string) error {
	glog.Infof("开始删除文件或文件夹，路径: %s", path)

	// 规范化路径
	path = filepath.Clean(path)
	if path == "." {
		path = ""
	}

	// 构建完整的路径
	fullPath := filepath.Join(consts.UploadDir, path)
	glog.Infof("目标路径: %s", fullPath)

	// 检查路径是否存在
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			glog.Infof("路径不存在，无需删除: %s", fullPath)
			return nil
		}
		glog.Errorf("检查路径状态失败: %s", err)
		return fmt.Errorf("检查路径状态失败: %s", err)
	}

	// 根据类型选择删除方法
	if info.IsDir() {
		// 检查目录是否为空
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			glog.Errorf("读取目录失败: %s", err)
			return fmt.Errorf("读取目录失败: %s", err)
		}
		if len(entries) > 0 {
			glog.Errorf("目录不为空: %s", fullPath)
			return fmt.Errorf("目录不为空，无法删除: %s", path)
		}
		err = os.Remove(fullPath) // 使用 Remove 删除空目录
	} else {
		err = os.Remove(fullPath) // 删除文件
	}

	if err != nil {
		glog.Errorf("删除失败: %s", err)
		return fmt.Errorf("删除失败: %s", err)
	}

	glog.Infof("删除成功: %s", fullPath)
	return nil
}

// UploadFile 上传文件
func (h *FileServiceImpl) UploadFile(path, fileName string, totalChunks int) error {
	glog.Infof("开始上传文件，路径: %s, 文件名: %s", path, fileName)

	// 规范化路径
	path = filepath.Clean(path)
	if path == "." {
		path = ""
	}

	// 构建完整的文件路径
	outFilePath := filepath.Join(consts.UploadDir, path, fileName)
	glog.Infof("目标文件路径: %s", outFilePath)

	// 创建目标目录
	targetDir := filepath.Dir(outFilePath)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		glog.Errorf("创建目标目录失败: %s, 路径: %s", err, targetDir)
		return fmt.Errorf("创建目标目录失败: %s", err)
	}

	// 检查目标目录是否有写入权限
	if err := h.checkDirectoryWritePermission(targetDir); err != nil {
		glog.Errorf("目标目录无写入权限: %s, 路径: %s", err, targetDir)
		return fmt.Errorf("目标目录无写入权限: %s", err)
	}

	// 检查文件是否已存在
	if _, err := os.Stat(outFilePath); err == nil {
		glog.Errorf("文件已存在: %s", outFilePath)
		return fmt.Errorf("文件已存在: %s", fileName)
	} else if !os.IsNotExist(err) {
		glog.Errorf("检查文件状态失败: %s", err)
		return fmt.Errorf("检查文件状态失败: %s", err)
	}

	glog.Infof("文件上传前置检查完成")
	return nil
}

// checkDirectoryWritePermission 检查目录是否有写入权限
func (h *FileServiceImpl) checkDirectoryWritePermission(dir string) error {
	// 创建临时文件
	tempFile := filepath.Join(dir, ".write_test")
	f, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	f.Close()

	// 删除临时文件
	return os.Remove(tempFile)
}

// GetFileStats 获取文件统计信息
func (s *FileServiceImpl) GetFileStats(path string) (*model.FileStats, error) {
	glog.Infof("开始获取文件统计信息，路径: %s", path)

	// 确保路径存在
	absPath := filepath.Join(consts.UploadDir, path)
	glog.Infof("完整路径: %s", absPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		glog.Errorf("路径不存在: %s", absPath)
		return nil, fmt.Errorf("路径不存在: %s", path)
	}

	stats := &model.FileStats{}

	// 遍历目录
	err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path != absPath { // 不统计根目录
				stats.TotalFolders++
			}
		} else {
			stats.TotalFiles++
			stats.TotalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		glog.Errorf("统计文件信息失败: %s", err)
		return nil, fmt.Errorf("统计文件信息失败: %s", err)
	}

	glog.Infof("成功获取文件统计信息: %+v", stats)
	return stats, nil
}

// SearchFiles 搜索文件
func (s *FileServiceImpl) SearchFiles(keyword string) ([]model.FileInfo, error) {
	glog.Infof("开始搜索文件，词: %s", keyword)

	var results []model.FileInfo
	err := filepath.Walk(consts.UploadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果文件名包含关键词
		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(keyword)) {
			relativePath, err := filepath.Rel(consts.UploadDir, path)
			if err != nil {
				glog.Errorf("获取相对路径失败: %s", err)
				return nil
			}

			fileInfo := model.FileInfo{
				FileName: info.Name(),
				FilePath: relativePath,
				FileSize: info.Size(),
				FileType: filepath.Ext(info.Name()),
				IsDir:    info.IsDir(),
				ModTime:  info.ModTime().Format(time.DateTime),
			}
			results = append(results, fileInfo)
		}
		return nil
	})

	if err != nil {
		glog.Errorf("搜索文件失败: %s", err)
		return nil, fmt.Errorf("搜索文件失败: %s", err)
	}

	glog.Infof("搜索完成，找到 %d 个匹配文件", len(results))
	return results, nil
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
		FilePath:  filePath,
		FileName:  info.Name(),
		IsDir:     info.IsDir(),
		CreatedAt: time.Now(),
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
func (h *FileServiceImpl) RemoveFile(path string) error {
	return h.DeleteFile(path)
}
