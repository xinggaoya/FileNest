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
	file.DELETE("/delete", fileController.DeleteFile)
	file.DELETE("/favorite", fileController.RemoveFavorite)
	file.POST("/rename", fileController.RenameFile)
	file.POST("/copy", fileController.CopyFile)
	file.POST("/move", fileController.MoveFile)
}

// RegisterGlobalMiddleware 注册全局中间件
func RegisterGlobalMiddleware(app *gin.Engine) {
	app.Use(cors.Default())

	app.Use(ginzap.Ginzap(glog.GetLogger(), time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	app.Use(ginzap.RecoveryWithZap(glog.GetLogger(), true))
}
