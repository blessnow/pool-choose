package main

import (
	"log"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/yuchi/cycle-stock/internal/handlers"
	"github.com/yuchi/cycle-stock/internal/repository"
)

func main() {
	// 初始化数据库
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/cycle_stock.db"
	}

	// 确保数据目录存在
	os.MkdirAll("./data", 0755)

	if err := repository.InitDB(dbPath); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 设置密码
	password := os.Getenv("PASSWORD")
	if password != "" {
		handlers.Password = password
	}

	// 创建路由
	r := gin.Default()

	// Gzip压缩（对静态文件和API响应都生效）
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// CORS配置
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API路由
	api := r.Group("/api")
	{
		// 认证
		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/logout", handlers.Logout)

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(handlers.AuthMiddleware)
		{
			// 股票
			auth.GET("/stocks", handlers.GetStocks)
			auth.GET("/stocks/:code", handlers.GetStock)
			// 注意：search/info 放到独立命名空间，避免与 /stocks/:code 路径冲突
			// （gin 在不同版本/平台对静态 vs 参数路由优先级处理可能不一致）
			auth.GET("/stock-search", handlers.SearchStocks)
			auth.GET("/stock-info", handlers.GetStockInfo)
			auth.POST("/stocks", handlers.CreateStock)
			auth.PUT("/stocks/:code", handlers.UpdateStock)
			auth.DELETE("/stocks/:code", handlers.DeleteStock)

			// 行情
			auth.GET("/quotes", handlers.GetQuotes)
			auth.GET("/company-summaries", handlers.GetCompanySummaries)
			auth.GET("/chart/:code", handlers.GetChart)
			auth.GET("/valuation/:code", handlers.GetValuation)
			auth.GET("/valuation-band", handlers.GetValuationBand)

			// 周期分析
			auth.GET("/cycle-insight", handlers.GetCycleInsight)
			auth.PUT("/cycle-insight", handlers.UpdateCycleInsight)

			// 持仓记录
			auth.GET("/positions", handlers.GetPositions)
			auth.POST("/positions", handlers.CreatePosition)
			auth.PUT("/positions/:id", handlers.UpdatePosition)
			auth.DELETE("/positions/:id", handlers.DeletePosition)
		}
	}

	// 静态文件服务（前端）
	r.Static("/assets", "./frontend/assets")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/favicon.ico", "./frontend/favicon.ico")

	// 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务启动在端口 %s", port)
	r.Run(":" + port)
}