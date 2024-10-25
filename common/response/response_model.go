package response

import "github.com/gin-gonic/gin"

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

func NewResponseModel(ctx *gin.Context, code int, message string, data interface{}) {
	model := &Res{
		Code:    code,
		Message: message,
		Data:    data,
	}
	ctx.JSON(200, model)
}

func NewPageResponseModel(ctx *gin.Context, code int, message string, total int64, data interface{}) {
	model := &Res{
		Code:    code,
		Message: message,
		Data: &PageRes{
			Total: total,
			List:  data,
		},
	}
	ctx.JSON(200, model)
}

// Success 成功
func Success(ctx *gin.Context, data any) {
	NewResponseModel(ctx, 1000, "success", data)
}

// Error 失败
func Error(ctx *gin.Context, message string) {
	NewResponseModel(ctx, 1001, message, nil)
}

// PageSuccess 分页成功
func PageSuccess(ctx *gin.Context, total int64, data any) {
	NewPageResponseModel(ctx, 1000, "success", total, data)
}

// PageError 分页失败
func PageError(ctx *gin.Context, message string) {
	NewResponseModel(ctx, 1001, message, nil)
}
