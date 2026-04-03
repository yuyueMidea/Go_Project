package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 简易请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		if raw != "" {
			path = path + "?" + raw
		}

		color := "\033[32m" // green
		if status >= 400 && status < 500 {
			color = "\033[33m" // yellow
		} else if status >= 500 {
			color = "\033[31m" // red
		}
		reset := "\033[0m"

		fmt.Printf("%s[%d]%s %v | %-7s %s (%v)\n",
			color, status, reset, start.Format("2006/01/02 15:04:05"),
			method, path, latency)
	}
}
