package router

import (
	"github.com/gin-gonic/gin"
	"my-api/handlers"
	"my-api/middleware"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.Cors())

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 待办事项 路由
		todoHandler := handlers.NewTodoHandler()
		todoGroup := api.Group("/todos")
		{
			todoGroup.POST("", todoHandler.Create)
			todoGroup.GET("", todoHandler.List)
			todoGroup.GET("/:id", todoHandler.GetByID)
			todoGroup.PUT("/:id", todoHandler.Update)
			todoGroup.DELETE("/:id", todoHandler.Delete)
			todoGroup.POST("/batch-delete", todoHandler.BatchDelete)
		}

	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
