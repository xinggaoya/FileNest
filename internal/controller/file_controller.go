package controller

import (
	"FileNest/common/response"
	"FileNest/internal/service"
	"github.com/gofiber/fiber/v3"
	"strconv"
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
	err = service.NewFileService().
		UploadFile(file, path, fileName, intIndexChunk, intTotalChunks)
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	return response.Success(ctx, nil)
}
