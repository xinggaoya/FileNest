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

func (h *FileController) UploadFile(ctx fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return response.Error(ctx, err.Error())
	}
	indexChunk := ctx.FormValue("indexChunk")
	totalChunks := ctx.FormValue("totalChunks")

	intIndexChunk, _ := strconv.Atoi(indexChunk)
	intTotalChunks, _ := strconv.Atoi(totalChunks)
	return response.Success(ctx, service.NewFileService().
		UploadFile(file, intIndexChunk, intTotalChunks))
}
