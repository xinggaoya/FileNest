package controller

import (
	"FileNest/common/glog"
	"FileNest/common/response"
	"FileNest/internal/consts"
	"FileNest/internal/service"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

/**
  @author: XingGao
  @date: 2024/9/28
**/

type FileController struct{}

func NewFileController() *FileController {
	return &FileController{}
}

// GetFileList 获取文件列表
func (h *FileController) GetFileList(ctx *gin.Context) {
	path := ctx.Query("path")

	list, err := service.NewFileService().GetFileList(path)
	if err != nil {
		response.Error(ctx, err.Error())
	}
	response.Success(ctx, list)
}

// DownloadFile 下载文件
func (h *FileController) DownloadFile(ctx *gin.Context) {
	path := ctx.Query("path")
	absPath, err := service.NewFileService().DownloadFile(path)
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

// CreateFolder 创建文件夹、
func (h *FileController) CreateFolder(ctx *gin.Context) {
	// 获取json path
	path := ctx.Query("path")
	if path == "" {
		response.Error(ctx, "path is empty")
		return
	}
	err := service.NewFileService().CreateDir(path)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

// DeleteFile 删除
func (h *FileController) DeleteFile(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		response.Error(ctx, "path is empty")
		return
	}
	err := service.NewFileService().DeleteFile(path)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

func (h *FileController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}
	indexChunk := ctx.PostForm("indexChunk")
	totalChunks := ctx.PostForm("totalChunks")
	fileName := ctx.PostForm("fileName")
	path := ctx.PostForm("path")
	// 是否覆盖上传
	override := ctx.PostForm("override")
	if override == "" {
		override = "false"
	}
	if fileName == "" || indexChunk == "" || totalChunks == "" {
		response.Error(ctx, "参数错误")
		return
	}

	intIndexChunk, _ := strconv.Atoi(indexChunk)
	intTotalChunks, _ := strconv.Atoi(totalChunks)
	overrideBool, _ := strconv.ParseBool(override)
	// 跳过重复上传
	if intIndexChunk == intTotalChunks && !overrideBool {
		// 检查文件是否存在
		filePath := filepath.Join(consts.UploadDir, path, fileName)
		if _, err = os.Stat(filePath); err == nil {
			response.NewResponseModel(ctx, 1005, "文件已存在", nil)
			return
		}
	}

	err = os.MkdirAll(consts.TempDir, os.ModePerm)
	if err != nil {
		glog.Errorf("unable to create temp directory: %s", err)
		return
	}
	name := filepath.Base(fileName)
	tempFileName := name + "_chunk_" + strconv.Itoa(intIndexChunk)
	chunkPath := filepath.Join(consts.TempDir, tempFileName)

	if err = ctx.SaveUploadedFile(file, chunkPath); err != nil {
		response.Error(ctx, err.Error())
		return
	}

	if intIndexChunk == intTotalChunks {
		err = service.NewFileService().
			UploadFile(path, fileName, intTotalChunks)
		if err != nil {
			response.Error(ctx, err.Error())
			return
		}

		// 删除临时文件
		go func() {
			time.Sleep(5 * time.Second)
			e := os.RemoveAll(consts.TempDir)
			if e != nil {
				glog.Errorf("delete temp file error: %s", e)
			}
			glog.Infof("delete temp file success")
		}()
	}
	response.Success(ctx, nil)
}
