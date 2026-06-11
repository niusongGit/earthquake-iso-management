package main

import (
	"earthquake-iso-management/internal/config"
	"earthquake-iso-management/internal/database"
	"earthquake-iso-management/internal/logger"
	"earthquake-iso-management/internal/router"
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed staticweb
var staticFiles embed.FS

func main() {
	// 加载配置
	if err := config.Load("config.json"); err != nil {
		fmt.Println("加载配置文件失败:", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init("logs"); err != nil {
		fmt.Println("初始化日志失败:", err)
		os.Exit(1)
	}
	defer logger.Log.Sync()

	// 初始化数据库
	if err := database.Init(); err != nil {
		logger.Log.Fatal("初始化数据库失败", zap.Error(err))
	}

	// 创建Gin引擎
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS中间件
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

	// 上传文件大小限制 50MB
	r.MaxMultipartMemory = 50 << 20

	// 构建静态文件子目录
	staticFS, err := fs.Sub(staticFiles, "staticweb")
	if err != nil {
		logger.Log.Fatal("读取静态文件目录失败", zap.Error(err))
	}

	// 注册路由
	router.Setup(r, staticFS)

	// 启动服务
	addr := fmt.Sprintf(":%d", config.C.Server.Port)
	logger.Log.Info("服务启动", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Log.Fatal("服务启动失败", zap.Error(err))
	}
}
