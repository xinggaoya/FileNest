package cmd

import (
	"FileNest/common/glog"
	"FileNest/router"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

/**
  @author: XingGao
  @date: 2024/9/22
**/

// Run 运行服务
func Run() {
	glog.Install()
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()
	router.Install(app)

	// 上传大小
	port := flag.Int("port", 9040, "port")
	flag.Parse()

	glog.Infof("starting server on port %d", *port)
	log.Fatal(app.Run(fmt.Sprintf(":%d", *port)))
}
