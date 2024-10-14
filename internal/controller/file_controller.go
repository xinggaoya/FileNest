package controller

import (
	"FileNest/common/glog"
	"FileNest/common/response"
	"FileNest/internal/consts"
	"FileNest/internal/service"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"io"
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
func (h *FileController) GetFileList(ctx fiber.Ctx) error {
	path := ctx.Query("path")

	list, err := service.NewFileService().GetFileList(path)
	if err != nil {
		return err
	}
	return response.Success(ctx, list)
}

// DownloadFile 下载文件
func (h *FileController) DownloadFile(ctx fiber.Ctx) error {
	path := ctx.Query("path")
	absPath, err := service.NewFileService().DownloadFile(path)
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	// 对文件名进行 URL 编码，并设置 Content-Disposition 响应头
	fileName := filepath.Base(absPath)
	encodedFileName := url.QueryEscape(fileName)
	ctx.Set("Content-Type", "application/octet-stream")
	ctx.Set("Content-Disposition", "attachment; filename*=UTF-8''"+encodedFileName)

	return ctx.SendFile(absPath)
}

// CreateFolder 创建文件夹、
func (h *FileController) CreateFolder(ctx fiber.Ctx) error {
	// 获取json path
	type Query struct {
		Path string `json:"path"`
	}
	var query Query
	err := ctx.Bind().Body(&query)
	if err != nil {
		return response.Error(ctx, "json error")
	}
	if query.Path == "" {
		return response.Error(ctx, "path is empty")
	}
	err = service.NewFileService().CreateDir(query.Path)
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	return response.Success(ctx, nil)
}

// DeleteFile 删除
func (h *FileController) DeleteFile(ctx fiber.Ctx) error {
	path := ctx.Query("path")
	if path == "" {
		return response.Error(ctx, "path is empty")
	}
	err := service.NewFileService().DeleteFile(path)
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	return response.Success(ctx, nil)
}

func (h *FileController) UploadFile(ctx fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	indexChunk := ctx.FormValue("indexChunk")
	totalChunks := ctx.FormValue("totalChunks")
	fileName := ctx.FormValue("fileName")
	path := ctx.FormValue("path")
	if fileName == "" || indexChunk == "" || totalChunks == "" {
		return response.Error(ctx, "参数错误")
	}

	intIndexChunk, _ := strconv.Atoi(indexChunk)
	intTotalChunks, _ := strconv.Atoi(totalChunks)

	err = os.MkdirAll(consts.TempDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create temp directory: %s", err)
	}
	name := filepath.Base(fileName)
	tempFileName := name + "_chunk_" + strconv.Itoa(intIndexChunk)
	chunkPath := filepath.Join(consts.TempDir, tempFileName)

	outFile, err := os.Create(chunkPath)
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	defer outFile.Close()

	f, err := file.Open()
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	defer f.Close()
	if _, err = io.Copy(outFile, f); err != nil {
		return response.Error(ctx, err.Error())
	}

	if intIndexChunk == intTotalChunks {
		err = service.NewFileService().
			UploadFile(path, fileName, intTotalChunks)
		if err != nil {
			return response.Error(ctx, err.Error())
		}

		// 删除临时文件
		go func() {
			time.Sleep(5 * time.Second)
			err = os.RemoveAll(consts.TempDir)
			if err != nil {
				glog.Errorf("unable to remove temp directory %s", err)
				return
			}
			glog.Info("清理临时文件成功")
		}()
	}
	return response.Success(ctx, nil)
}
