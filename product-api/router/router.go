package router

import (
	"github.com/gin-gonic/gin"
	"generated-api/handlers"
	"generated-api/middleware"
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
		// 商品 路由
		productHandler := handlers.NewProductHandler()
		productGroup := api.Group("/products")
		{
			productGroup.POST("", productHandler.Create)
			productGroup.GET("", productHandler.List)
			productGroup.GET("/:id", productHandler.GetByID)
			productGroup.PUT("/:id", productHandler.Update)
			productGroup.DELETE("/:id", productHandler.Delete)
			productGroup.POST("/batch-delete", productHandler.BatchDelete)
		}

	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
