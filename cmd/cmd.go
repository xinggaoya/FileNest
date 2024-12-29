package cmd

import (
	"FileNest/common/glog"
	"FileNest/internal/cache"
	"FileNest/router"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

/**
  @author: XingGao
  @date: 2024/9/22
**/

// Run 运行服务
func Run() {
	// 初始化日志
	glog.Install()

	// 初始化 Redis
	if err := cache.InitRedis(); err != nil {
		glog.Errorf("初始化 Redis 失败: %s", err)
		os.Exit(1)
	}

	// 程序退出时清理资源
	defer func() {
		if err := cache.Close(); err != nil {
			glog.Errorf("关闭 Redis 连接失败: %s", err)
		}
	}()

	// 监听系统信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		glog.Infof("收到系统信号: %v, 开始清理资源...", sig)
		if err := cache.Close(); err != nil {
			glog.Errorf("关闭 Redis 连接失败: %s", err)
		}
		os.Exit(0)
	}()

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	router.Install(app)

	// 上传大小
	port := flag.Int("port", 9040, "port")
	flag.Parse()

	glog.Infof("starting server on port %d", *port)
	log.Fatal(app.Run(fmt.Sprintf(":%d", *port)))
}
