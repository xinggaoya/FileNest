package controller

import (
	"FileNest/common/glog"
	"FileNest/internal/service"
	"FileNest/internal/utils/response"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShareController 分享控制器
type ShareController struct {
	shareService service.ShareService
}

// NewShareController 创建分享控制器实例
func NewShareController(shareService service.ShareService) *ShareController {
	return &ShareController{
		shareService: shareService,
	}
}

// CreateShare 创建分享
func (h *ShareController) CreateShare(ctx *gin.Context) {
	var req struct {
		FilePath    string `json:"filePath" binding:"required"`
		Password    string `json:"password"`
		ExpireHours int    `json:"expireHours"`
		MaxDownload int    `json:"maxDownload"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		glog.Errorf("解析请求参数失败: %s", err)
		response.Error(ctx, "请求参数格式错误")
		return
	}

	clientIP := ctx.ClientIP()
	glog.Infof("收到创建分享请求，文件路径: %s, 客户端IP: %s", req.FilePath, clientIP)

	share, err := h.shareService.CreateShare(
		req.FilePath,
		req.Password,
		req.ExpireHours,
		req.MaxDownload,
		req.Description,
		clientIP,
	)

	if err != nil {
		glog.Errorf("创建分享失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("分享创建成功，分享码: %s", share.ShareCode)
	response.Success(ctx, share)
}

// GetShareInfo 获取分享信息
func (h *ShareController) GetShareInfo(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	if shareCode == "" {
		response.Error(ctx, "分享码不能为空")
		return
	}

	glog.Infof("收到获取分享信息请求，分享码: %s", shareCode)

	share, err := h.shareService.GetShareByCode(shareCode)
	if err != nil {
		glog.Errorf("获取分享信息失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	// 隐藏敏感信息
	shareInfo := map[string]interface{}{
		"shareCode":   share.ShareCode,
		"fileName":    share.FileName,
		"fileSize":    share.FileSize,
		"isDir":       share.IsDir,
		"hasPassword": share.Password != "",
		"expireTime":  share.ExpireTime,
		"maxDownload": share.MaxDownload,
		"downloaded":  share.Downloaded,
		"description": share.Description,
		"createTime":  share.CreateTime,
	}

	response.Success(ctx, shareInfo)
}

// ValidateShare 验证分享
func (h *ShareController) ValidateShare(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	password := ctx.PostForm("password")

	glog.Infof("收到验证分享请求，分享码: %s", shareCode)

	share, err := h.shareService.ValidateShare(shareCode, password)
	if err != nil {
		glog.Errorf("验证分享失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	// 返回验证成功信息
	result := map[string]interface{}{
		"valid":    true,
		"fileName": share.FileName,
		"fileSize": share.FileSize,
		"isDir":    share.IsDir,
	}

	response.Success(ctx, result)
}

// DownloadSharedFile 下载分享文件
func (h *ShareController) DownloadSharedFile(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	password := ctx.Query("password")

	glog.Infof("收到下载分享文件请求，分享码: %s", shareCode)

	absPath, err := h.shareService.DownloadSharedFile(shareCode, password)
	if err != nil {
		glog.Errorf("下载分享文件失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	// 设置下载响应头
	fileName := filepath.Base(absPath)
	encodedFileName := url.QueryEscape(fileName)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFileName)

	glog.Infof("开始下载分享文件: %s", fileName)
	ctx.File(absPath)
}

// GetMyShares 获取我的分享列表
func (h *ShareController) GetMyShares(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	clientIP := ctx.ClientIP()
	glog.Infof("收到获取我的分享列表请求，客户端IP: %s, 页码: %d, 页大小: %d", clientIP, page, pageSize)

	shares, total, err := h.shareService.GetMyShares(clientIP, page, pageSize)
	if err != nil {
		glog.Errorf("获取分享列表失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	result := map[string]interface{}{
		"shares":   shares,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}

	glog.Infof("获取分享列表成功，共 %d 条记录", total)
	response.Success(ctx, result)
}

// DeleteShare 删除分享
func (h *ShareController) DeleteShare(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	if shareCode == "" {
		response.Error(ctx, "分享码不能为空")
		return
	}

	glog.Infof("收到删除分享请求，分享码: %s", shareCode)

	// 验证是否是分享创建者（通过IP验证）
	clientIP := ctx.ClientIP()
	share, err := h.shareService.GetShareByCode(shareCode)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	if share.CreatorIP != clientIP {
		glog.Errorf("删除分享失败：无权限，分享创建者IP: %s, 当前IP: %s", share.CreatorIP, clientIP)
		response.Error(ctx, "无权限删除此分享")
		return
	}

	if err := h.shareService.DeleteShare(shareCode); err != nil {
		glog.Errorf("删除分享失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("删除分享成功: %s", shareCode)
	response.Success(ctx, nil)
}

// DisableShare 禁用分享
func (h *ShareController) DisableShare(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	if shareCode == "" {
		response.Error(ctx, "分享码不能为空")
		return
	}

	glog.Infof("收到禁用分享请求，分享码: %s", shareCode)

	// 验证是否是分享创建者（通过IP验证）
	clientIP := ctx.ClientIP()
	share, err := h.shareService.GetShareByCode(shareCode)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	if share.CreatorIP != clientIP {
		glog.Errorf("禁用分享失败：无权限，分享创建者IP: %s, 当前IP: %s", share.CreatorIP, clientIP)
		response.Error(ctx, "无权限禁用此分享")
		return
	}

	if err := h.shareService.DisableShare(shareCode); err != nil {
		glog.Errorf("禁用分享失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	glog.Infof("禁用分享成功: %s", shareCode)
	response.Success(ctx, nil)
}

// GetShareStats 获取分享统计信息
func (h *ShareController) GetShareStats(ctx *gin.Context) {
	shareCode := ctx.Param("shareCode")
	if shareCode == "" {
		response.Error(ctx, "分享码不能为空")
		return
	}

	glog.Infof("收到获取分享统计请求，分享码: %s", shareCode)

	share, err := h.shareService.GetShareStats(shareCode)
	if err != nil {
		glog.Errorf("获取分享统计失败: %s", err)
		response.Error(ctx, err.Error())
		return
	}

	// 验证是否是分享创建者（通过IP验证）
	clientIP := ctx.ClientIP()
	if share.CreatorIP != clientIP {
		response.Error(ctx, "无权限查看此分享统计")
		return
	}

	response.Success(ctx, share)
}
