package cmd

import (
	"FileNest/common/glog"
	"FileNest/common/response"
	"FileNest/router"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

/**
  @author: XingGao
  @date: 2024/9/22
**/

// 配置 Fiber 应用实例
func newFiberConfig() fiber.Config {
	return fiber.Config{
		BodyLimit:                    1024 * 1024 * 1024 * 10,
		DisablePreParseMultipartForm: true,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return response.Error(ctx, "系统出现异常")
		},
		StreamRequestBody: true, // 开启流式处理
	}
}

// Run 运行服务
func Run() {
	glog.Install()

	app := fiber.New(newFiberConfig())
	router.Install(app)

	// 上传大小
	port := flag.Int("port", 9040, "port")
	flag.Parse()
	log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
}
