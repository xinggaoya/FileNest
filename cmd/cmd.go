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
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			return response.Error(ctx, "系统出现异常")
		},
	}
}

// Run 运行服务
func Run() {
	glog.Install()

	app := fiber.New(newFiberConfig())
	router.Install(app)

	port := flag.Int("port", 3000, "port")
	flag.Parse()
	log.Fatal(app.Listen(fmt.Sprintf(":%d", *port)))
}
