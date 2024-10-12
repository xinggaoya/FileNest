package response

import (
	"github.com/gofiber/fiber/v3"
)

/**
  @author: XingGao
  @date: 2024/4/19
**/

type Res struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageRes struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

func NewResponseModel(ctx fiber.Ctx, code int, message string, data interface{}) error {
	model := &Res{
		Code:    code,
		Message: message,
		Data:    data,
	}
	err := ctx.Status(200).JSON(model)
	if err != nil {
		return err
	}
	return nil
}

func NewPageResponseModel(ctx fiber.Ctx, code int, message string, total int64, data interface{}) error {
	model := &Res{
		Code:    code,
		Message: message,
		Data: &PageRes{
			Total: total,
			List:  data,
		},
	}
	err := ctx.Status(200).JSON(model)
	if err != nil {
		return err
	}
	return nil
}

// Success 成功
func Success(ctx fiber.Ctx, data any) error {
	return NewResponseModel(ctx, 1000, "success", data)
}

// Error 失败
func Error(ctx fiber.Ctx, message string) error {
	return NewResponseModel(ctx, 1001, message, nil)
}

// PageSuccess 分页成功
func PageSuccess(ctx fiber.Ctx, total int64, data any) error {
	return NewPageResponseModel(ctx, 1000, "success", total, data)
}

// PageError 分页失败
func PageError(ctx fiber.Ctx, message string) error {
	return NewResponseModel(ctx, 1001, message, nil)
}
