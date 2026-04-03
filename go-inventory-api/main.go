package main

import (
	"log"
	"net/http"
	"os"

	"inventory-api/handlers"
	"inventory-api/middleware"
	"inventory-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// ── 1. 初始化数据库 ───────────────────────────────────────────────────────
	dsn := getEnv("DB_PATH", "inventory.db")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 启用 WAL 模式（提升并发性能）
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA foreign_keys=ON;")

	// 自动建表
	if err := db.AutoMigrate(&models.Product{}, &models.StockLog{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库初始化完成")

	// ── 2. 初始化路由 ─────────────────────────────────────────────────────────
	gin.SetMode(getEnv("GIN_MODE", gin.DebugMode))
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.OK(gin.H{"status": "ok"}))
	})

	// ── 3. 注册路由 ───────────────────────────────────────────────────────────
	productH := handlers.NewProductHandler(db)
	stockH := handlers.NewStockHandler(db)

	api := r.Group("/api/v1")
	{
		// 商品管理
		products := api.Group("/products")
		{
			products.POST("", productH.CreateProduct)    // 新增商品
			products.GET("", productH.ListProducts)      // 商品列表（支持分页/搜索）
			products.GET("/:id", productH.GetProduct)    // 商品详情
			products.PUT("/:id", productH.UpdateProduct) // 更新商品
			products.DELETE("/:id", productH.DeleteProduct) // 删除商品（软删除）
		}

		// 入库 / 出库 / 流水
		stock := api.Group("/stock")
		{
			stock.POST("/in", stockH.StockIn)                        // 入库
			stock.POST("/out", stockH.StockOut)                      // 出库
			stock.GET("/logs", stockH.ListStockLogs)                 // 流水查询
			stock.GET("/:product_id", stockH.GetProductStock)        // 单品库存
		}
	}

	// ── 4. 启动服务 ───────────────────────────────────────────────────────────
	addr := getEnv("ADDR", ":8080")
	log.Printf("🚀 服务启动: http://localhost%s", addr)
	printRoutes()

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func printRoutes() {
	log.Println("──────────────────────────────────────────────")
	log.Println("  商品管理")
	log.Println("  POST   /api/v1/products          新增商品")
	log.Println("  GET    /api/v1/products           商品列表")
	log.Println("  GET    /api/v1/products/:id       商品详情")
	log.Println("  PUT    /api/v1/products/:id       更新商品")
	log.Println("  DELETE /api/v1/products/:id       删除商品")
	log.Println("  入库/出库")
	log.Println("  POST   /api/v1/stock/in           入库")
	log.Println("  POST   /api/v1/stock/out          出库")
	log.Println("  GET    /api/v1/stock/logs         流水查询")
	log.Println("  GET    /api/v1/stock/:product_id  单品库存")
	log.Println("──────────────────────────────────────────────")
}
