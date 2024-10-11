package middlewares

import (
	"FileNest/common/glog"
	"github.com/gofiber/fiber/v3"
	"time"
)

// Logger 记录API日志
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 以当前时间计耗时
		start := time.Now()
		err := c.Next()
		if err != nil {
			return err
		}
		glog.Infof(
			"%s | %s | %dms | %d | %s | %s",
			c.IP(),
			c.Method(),
			time.Since(start).Milliseconds(),
			c.Response().StatusCode(),
			c.Path(),
			c.Protocol(),
		)
		return nil
	}
}
