package router

import (
	"FileNest/common/glog"
	"FileNest/internal/controller"
	"FileNest/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

/**
  @author: XingGao
  @date: 2024/9/22
**/

// Install 安装路由
func Install(app *gin.Engine) {

	RegisterGlobalMiddleware(app)

	index := app.Group("/")

	fileController := controller.NewFileController(service.NewFileService())
	shareController := controller.NewShareController(service.NewShareService())

	api := index.Group("/api")

	file := api.Group("/file")
	file.GET("/list", fileController.GetFileList)
	file.GET("/stats", fileController.GetFileStats)
	file.GET("/search", fileController.SearchFiles)
	file.GET("/favorites", fileController.GetFavorites)
	file.POST("/create-folder", fileController.CreateFolder)
	file.POST("/upload", fileController.UploadFile)
	file.POST("/upload-chunk", fileController.UploadChunk)
	file.POST("/merge-chunks", fileController.MergeChunks)
	file.POST("/favorite", fileController.AddFavorite)
	file.GET("/download", fileController.DownloadFile)
	file.GET("/preview", fileController.PreviewFile)
	file.DELETE("/delete", fileController.DeleteFile)
	file.DELETE("/favorite", fileController.RemoveFavorite)
	file.POST("/rename", fileController.RenameFile)
	file.POST("/copy", fileController.CopyFile)
	file.POST("/move", fileController.MoveFile)

	// 分享相关路由
	share := api.Group("/share")
	share.POST("/create", shareController.CreateShare)                    // 创建分享
	share.GET("/info/:shareCode", shareController.GetShareInfo)           // 获取分享信息
	share.POST("/validate/:shareCode", shareController.ValidateShare)     // 验证分享
	share.GET("/download/:shareCode", shareController.DownloadSharedFile) // 下载分享文件
	share.GET("/my", shareController.GetMyShares)                         // 获取我的分享列表
	share.DELETE("/:shareCode", shareController.DeleteShare)              // 删除分享
	share.PUT("/disable/:shareCode", shareController.DisableShare)        // 禁用分享
	share.GET("/stats/:shareCode", shareController.GetShareStats)         // 获取分享统计
}

// RegisterGlobalMiddleware 注册全局中间件
func RegisterGlobalMiddleware(app *gin.Engine) {
	app.Use(cors.Default())

	app.Use(ginzap.Ginzap(glog.GetLogger(), time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	app.Use(ginzap.RecoveryWithZap(glog.GetLogger(), true))
}
