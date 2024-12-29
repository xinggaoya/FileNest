package controller

import (
	"FileNest/common/glog"
	"FileNest/internal/consts"
	"FileNest/internal/service"
	"FileNest/internal/utils/response"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService service.FileService
}

func NewFileController(fileService service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// GetFileList 获取文件列表
func (h *FileController) GetFileList(ctx *gin.Context) {
	path := ctx.Query("path")
	glog.Infof("收到获取文件列表请求，路径: %s", path)

	list, err := h.fileService.GetFileList(path)
	if err != nil {
		glog.Errorf("获取文件列表失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}
	glog.Infof("成功获取文件列表，返回 %d 个文件", len(list))
	response.Success(ctx, list)
}

// DownloadFile 下载文件
func (h *FileController) DownloadFile(ctx *gin.Context) {
	path := ctx.Query("path")
	absPath, err := h.fileService.DownloadFile(path)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	// 对文件名进行 URL 编码，并设置 Content-Disposition 响应头
	fileName := filepath.Base(absPath)
	encodedFileName := url.QueryEscape(fileName)
	ctx.Set("Content-Type", "application/octet-stream")
	ctx.Set("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFileName)

	ctx.File(absPath)
}

// CreateFolder 创建文件夹
func (h *FileController) CreateFolder(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		response.Error(ctx, "path is empty")
		return
	}
	err := h.fileService.CreateFolder(path)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

// DeleteFile 删除文件
func (h *FileController) DeleteFile(ctx *gin.Context) {
	path := ctx.Query("path")
	force := ctx.Query("force") == "true"
	glog.Infof("收到删除文件请求，路径: %s, 强制删除: %v", path, force)

	err := h.fileService.DeleteFile(path, force)
	if err != nil {
		glog.Errorf("删除文件失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("成功删除文件: %s", path)
	response.Success(ctx, nil)
}

// UploadFile 上传文件
func (h *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		glog.Errorf("获取上传文件失败: %s", err)
		response.Error(ctx, "获取上传文件失败")
		return
	}

	fileName := ctx.PostForm("fileName")
	if fileName == "" {
		fileName = file.Filename
	}

	path := ctx.PostForm("path")
	override := ctx.PostForm("override") == "true"
	glog.Infof("收到文件上传请求，文件名: %s, 路径: %s, 是否覆盖: %v", fileName, path, override)

	// 规范化路径
	path = filepath.Clean(path)
	if path == "." {
		path = ""
	}

	// 确保目标目录存在
	uploadPath := filepath.Join(consts.UploadDir, path)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		glog.Errorf("创建目标目录失败: %s, 路径: %s", err, uploadPath)
		response.Error(ctx, "创建目标目录失败")
		return
	}

	// 检查目标目录是否有写入权限
	if err := h.checkDirectoryWritePermission(uploadPath); err != nil {
		glog.Errorf("目标目录无写入权限: %s, 路径: %s", err, uploadPath)
		response.Error(ctx, "目标目录无写入权限")
		return
	}

	// 检查文件上传前置条件
	if err := h.fileService.UploadFile(path, fileName, 1, override); err != nil {
		glog.Errorf("文件上传前置检查失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadPath, fileName)
	if err = ctx.SaveUploadedFile(file, filePath); err != nil {
		glog.Errorf("保存文件失败: %s, 路径: %s", err, filePath)
		response.Error(ctx, "保存文件失败")
		return
	}

	// 设置文件权限
	if err = os.Chmod(filePath, 0644); err != nil {
		glog.Warnf("设置文件权限失败: %s, 路径: %s", err, filePath)
	}

	glog.Infof("文件上传成功: %s", filePath)
	response.Success(ctx, map[string]string{
		"path": filepath.Join(path, fileName),
	})
}

// UploadChunk 上传文件分块
func (h *FileController) UploadChunk(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		glog.Errorf("获取上传文件失败: %s", err)
		response.Error(ctx, "获取上传文件失败")
		return
	}

	fileName := ctx.PostForm("fileName")
	if fileName == "" {
		fileName = file.Filename
	}

	path := ctx.PostForm("path")
	override := ctx.PostForm("override") == "true"
	chunkIndex := ctx.PostForm("chunkIndex")
	totalChunks := ctx.PostForm("totalChunks")

	// 转换chunkIndex和totalChunks为整数
	chunkIndexInt, err := strconv.Atoi(chunkIndex)
	if err != nil {
		glog.Errorf("分块索引格式错误: %s", err)
		response.Error(ctx, "分块索引格式错误")
		return
	}

	totalChunksInt, err := strconv.Atoi(totalChunks)
	if err != nil {
		glog.Errorf("总分块数格式错误: %s", err)
		response.Error(ctx, "总分块数格式错误")
		return
	}

	glog.Infof("收到文件分块上传请求，文件名: %s, 路径: %s, 是否覆盖: %v, 分块索引: %d, 总分块数: %d",
		fileName, path, override, chunkIndexInt, totalChunksInt)

	// 规范化路径
	path = filepath.Clean(path)
	if path == "." {
		path = ""
	}

	// 确保临时目录存在
	tempDir := filepath.Join(consts.TempDir, path, fileName)
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		glog.Errorf("创建临时目录失败: %s, 路径: %s", err, tempDir)
		response.Error(ctx, "创建临时目录失败")
		return
	}

	// 保存分块文件
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkIndexInt))
	if err = ctx.SaveUploadedFile(file, chunkPath); err != nil {
		glog.Errorf("保存分块文件失败: %s, 路径: %s", err, chunkPath)
		response.Error(ctx, "保存分块文件失败")
		return
	}

	response.Success(ctx, map[string]interface{}{
		"chunkIndex": chunkIndexInt,
		"path":       filepath.Join(path, fileName),
	})
}

// MergeChunks 合并文件分块
func (h *FileController) MergeChunks(ctx *gin.Context) {
	var req struct {
		FileName    string `json:"fileName"`
		Path        string `json:"path"`
		TotalChunks int    `json:"totalChunks"`
		Override    bool   `json:"override"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		glog.Errorf("解析请求参数失败: %s", err)
		response.Error(ctx, "解析请求参数失败")
		return
	}

	glog.Infof("收到合并文件请求，文件名: %s, 路径: %s, 总分块数: %d, 是否覆盖: %v",
		req.FileName, req.Path, req.TotalChunks, req.Override)

	// 规范化路径
	req.Path = filepath.Clean(req.Path)
	if req.Path == "." {
		req.Path = ""
	}

	// 检查文件上传前置条件
	if err := h.fileService.UploadFile(req.Path, req.FileName, req.TotalChunks, req.Override); err != nil {
		glog.Errorf("文件上传前置检查失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	// 确保目标目录存在
	uploadPath := filepath.Join(consts.UploadDir, req.Path)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		glog.Errorf("创建目标目录失败: %s, 路径: %s", err, uploadPath)
		response.Error(ctx, "创建目标目录失败")
		return
	}

	// 检查目标目录是否有写入权限
	if err := h.checkDirectoryWritePermission(uploadPath); err != nil {
		glog.Errorf("目标目录无写入权限: %s, 路径: %s", err, uploadPath)
		response.Error(ctx, "目标目录无写入权限")
		return
	}

	// 合并文件
	tempDir := filepath.Join(consts.TempDir, req.Path, req.FileName)
	targetFile := filepath.Join(uploadPath, req.FileName)

	// 创建目标文件
	outFile, err := os.Create(targetFile)
	if err != nil {
		glog.Errorf("创建目标文件失败: %s, 路径: %s", err, targetFile)
		response.Error(ctx, "创建目标文件失败")
		return
	}
	defer outFile.Close()

	// 按顺序合并分块
	for i := 0; i < req.TotalChunks; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", i))

		// 读取分块文件
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			glog.Errorf("打开分块文件失败: %s, 路径: %s", err, chunkPath)
			response.Error(ctx, "打开分块文件失败")
			return
		}

		// 复制分块内容到目标文件
		if _, err = io.Copy(outFile, chunkFile); err != nil {
			chunkFile.Close()
			glog.Errorf("复制分块内容失败: %s, 路径: %s", err, chunkPath)
			response.Error(ctx, "复制分块内容失败")
			return
		}
		chunkFile.Close()

		// 删除分块文件
		if err = os.Remove(chunkPath); err != nil {
			glog.Warnf("删除分块文件失败: %s, 路径: %s", err, chunkPath)
		}
	}

	// 清理临时目录
	if err := os.RemoveAll(tempDir); err != nil {
		glog.Warnf("清理临时目录失败: %s, 路径: %s", err, tempDir)
	}

	// 清除相关缓存
	h.fileService.ClearFileCache(req.Path)

	// 设置文件权限
	if err = os.Chmod(targetFile, 0644); err != nil {
		glog.Warnf("设置文件权限失败: %s, 路径: %s", err, targetFile)
	}

	glog.Infof("文件合并成功: %s", targetFile)
	response.Success(ctx, map[string]string{
		"path": filepath.Join(req.Path, req.FileName),
	})
}

// GetFileStats 获取文件统计信息
func (h *FileController) GetFileStats(ctx *gin.Context) {
	path := ctx.Query("path")
	glog.Infof("收到获取文件统计信息请求，路径: %s", path)

	stats, err := h.fileService.GetFileStats(path)
	if err != nil {
		glog.Errorf("获取文件统计信息失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("成功获取文件统计信息: %+v", stats)
	response.Success(ctx, stats)
}

// SearchFiles 搜索文件
func (h *FileController) SearchFiles(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	glog.Infof("收到搜索文件请求，关键词: %s", keyword)

	if keyword == "" {
		glog.Errorf("搜索关键词为空")
		response.Error(ctx, "搜索关键词不能为空")
		return
	}

	files, err := h.fileService.SearchFiles(keyword)
	if err != nil {
		glog.Errorf("搜索文件失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("搜索完成，找到 %d 个匹配文件", len(files))
	response.Success(ctx, files)
}

// AddFavorite 添加收藏
func (h *FileController) AddFavorite(ctx *gin.Context) {
	path := ctx.Query("path")
	glog.Infof("收到添加收藏请求，路径: %s", path)

	if path == "" {
		glog.Errorf("文件路径为空")
		response.Error(ctx, "文件路径不能为空")
		return
	}

	err := h.fileService.AddFavorite(path)
	if err != nil {
		glog.Errorf("添加收藏失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("添加收藏成功: %s", path)
	response.Success(ctx, nil)
}

// RemoveFavorite 取消收藏
func (h *FileController) RemoveFavorite(ctx *gin.Context) {
	path := ctx.Query("path")
	glog.Infof("收到取消收藏请求，路径: %s", path)

	if path == "" {
		glog.Errorf("文件路径为空")
		response.Error(ctx, "文件路径不能为空")
		return
	}

	err := h.fileService.RemoveFavorite(path)
	if err != nil {
		glog.Errorf("取消收藏失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("取消收藏成功: %s", path)
	response.Success(ctx, nil)
}

// GetFavorites 获取收藏列表
func (h *FileController) GetFavorites(ctx *gin.Context) {
	glog.Info("收到获取收藏列表请求")

	favorites, err := h.fileService.GetFavorites()
	if err != nil {
		glog.Errorf("获取收藏列表失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("获取收藏列表成功，共 %d 条记录", len(favorites))
	response.Success(ctx, favorites)
}

// RenameFile 重命名文件或文件夹
func (h *FileController) RenameFile(ctx *gin.Context) {
	oldPath := ctx.Query("path")
	newName := ctx.Query("newName")

	glog.Infof("收到重命名请求，原路径: %s, 新名称: %s", oldPath, newName)

	if oldPath == "" || newName == "" {
		glog.Errorf("参数错误：路径或新名称为空")
		response.Error(ctx, "路径和新名称不能为空")
		return
	}

	if err := h.fileService.RenameFile(oldPath, newName); err != nil {
		glog.Errorf("重命名失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Info("重命名成功")
	response.Success(ctx, nil)
}

// CopyFile 复制文件或文件夹
func (h *FileController) CopyFile(ctx *gin.Context) {
	srcPath := ctx.Query("srcPath")
	destPath := ctx.Query("destPath")

	glog.Infof("收到复制请求，源路径: %s, 目标路径: %s", srcPath, destPath)

	if srcPath == "" || destPath == "" {
		glog.Errorf("参数错误：源路径或目标路径为空")
		response.Error(ctx, "源路径和目标路径不能为空")
		return
	}

	if err := h.fileService.CopyFile(srcPath, destPath); err != nil {
		glog.Errorf("复制失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Info("复制成功")
	response.Success(ctx, nil)
}

// MoveFile 移动文件或文件夹
func (h *FileController) MoveFile(ctx *gin.Context) {
	srcPath := ctx.Query("srcPath")
	destPath := ctx.Query("destPath")

	glog.Infof("收到移动请求，源路径: %s, 目标路径: %s", srcPath, destPath)

	if srcPath == "" || destPath == "" {
		glog.Errorf("参数错误：源路径或目标路径为空")
		response.Error(ctx, "源路径和目标路径不能为空")
		return
	}

	if err := h.fileService.MoveFile(srcPath, destPath); err != nil {
		glog.Errorf("移动失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Info("移动成功")
	response.Success(ctx, nil)
}

// checkDirectoryWritePermission 检查目录是否有写入权限
func (h *FileController) checkDirectoryWritePermission(dir string) error {
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
