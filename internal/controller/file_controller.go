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
