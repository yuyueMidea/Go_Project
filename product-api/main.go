package main

import (
	"flag"
	"fmt"
	"log"

	"generated-api/database"
	"generated-api/router"
)

func main() {
	// 命令行参数
	port := flag.String("port", "8080", "服务端口")
	dbPath := flag.String("db", "data.db", "SQLite数据库文件路径")
	flag.Parse()

	// 初始化数据库
	if err := database.InitDB(*dbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 配置路由
	r := router.SetupRouter()

	// 启动服务
	addr := fmt.Sprintf(":%s", *port)
	log.Printf("🚀 服务启动成功，监听地址: http://localhost:%s", *port)
	log.Printf("📋 健康检查: http://localhost:%s/health", *port)
	log.Printf("📖 API基础路径: http://localhost:%s/api/v1", *port)
	log.Println("========================================")
	log.Println("  📁 商品: /api/v1/products")

	log.Println("========================================")

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
