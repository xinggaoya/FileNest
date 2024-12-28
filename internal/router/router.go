package router

import (
	"FileNest/internal/controller"
	"FileNest/internal/service"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 允许跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 设置静态文件目录
	r.Static("/upload", "./upload")

	// API 路由组
	api := r.Group("/api")
	{
		// 文件相关路由
		file := api.Group("/file")
		{
			fileController := controller.NewFileController(service.NewFileService())
			file.GET("/list", fileController.GetFileList)
			file.GET("/stats", fileController.GetFileStats)
			file.GET("/search", fileController.SearchFiles)
			file.GET("/favorites", fileController.GetFavorites)
			file.POST("/create-folder", fileController.CreateFolder)
			file.POST("/upload", fileController.UploadFile)
			file.POST("/favorite", fileController.AddFavorite)
			file.GET("/download", fileController.DownloadFile)
			file.DELETE("/delete", fileController.DeleteFile)
			file.DELETE("/favorite", fileController.RemoveFavorite)
		}
	}

	return r
}
