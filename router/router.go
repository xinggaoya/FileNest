package router

import (
	"FileNest/common/glog"
	"FileNest/internal/controller"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

/**
  @author: XingGao
  @date: 2024/9/22
**/

// Install 安装路由
func Install(app *gin.Engine) {

	RegisterGlobalMiddleware(app)

	index := app.Group("/")

	helloWordController := controller.NewFileController()

	api := index.Group("/api")

	file := api.Group("/file")
	file.GET("/list", helloWordController.GetFileList)
	file.POST("/create", helloWordController.CreateFolder)
	file.DELETE("/delete", helloWordController.DeleteFile)
	file.POST("/upload", helloWordController.UploadFile)
	file.GET("/download", helloWordController.DownloadFile)
}

// RegisterGlobalMiddleware 注册全局中间件
func RegisterGlobalMiddleware(app *gin.Engine) {
	app.Use(cors.Default())

	app.Use(ginzap.Ginzap(glog.GetLogger(), time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	app.Use(ginzap.RecoveryWithZap(glog.GetLogger(), true))
}
