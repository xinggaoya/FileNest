package router

import (
	"FileNest/common/middlewares"
	"FileNest/internal/controller"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

/**
  @author: XingGao
  @date: 2024/9/22
**/

// Install 安装路由
func Install(app *fiber.App) {

	RegisterGlobalMiddleware(app)

	index := app.Group("/")

	helloWordController := controller.NewFileController()

	api := index.Group("/api")

	file := api.Group("/file")
	file.Get("/list", helloWordController.GetFileList)
	file.Post("/create", helloWordController.CreateFolder)
	file.Delete("/delete", helloWordController.DeleteFile)
	file.Post("/upload", helloWordController.UploadFile)
	file.Get("/download", helloWordController.DownloadFile)
}

// RegisterGlobalMiddleware 注册全局中间件
func RegisterGlobalMiddleware(app *fiber.App) {
	// 跨域
	app.Use(cors.New())
	// 限流
	//app.Use(limiter.New())
	// 重复请求
	app.Use(idempotency.New())
	// 错误处理
	app.Use(recover.New())
	// 日志
	app.Use(middlewares.Logger())
	app.Use(pprof.New(pprof.Config{Prefix: "/endpoint-prefix"}))
}
