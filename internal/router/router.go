package router

import (
	"earthquake-iso-management/internal/handler"
	"earthquake-iso-management/internal/middleware"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, staticFS fs.FS) {
	// 前台API（无需认证）
	front := r.Group("/api/front")
	{
		front.GET("/documents", handler.GetFrontDocumentList)
		front.GET("/documents/:id", handler.GetFrontDocumentByID)
		front.GET("/documents/:id/preview", handler.PreviewAttachment)
		front.GET("/documents/:id/download", handler.DownloadAttachment)
	}

	// 后台API
	admin := r.Group("/api/admin")
	{
		admin.POST("/login", handler.AdminLogin)

		// 需要认证的接口
		auth := admin.Group("").Use(middleware.AuthRequired())
		{
			auth.POST("/documents", handler.CreateDocument)
			auth.GET("/documents", handler.GetAdminDocumentList)
			auth.GET("/documents/:id", handler.GetAdminDocumentByID)
			auth.PUT("/documents/:id", handler.UpdateDocument)
			auth.DELETE("/documents/:id", handler.DeleteDocument)
		}
	}

	// 静态文件服务
	staticServer := http.FileServer(http.FS(staticFS))
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 尝试读取静态文件
		if path == "/" {
			path = "/index.html"
		}

		// 检查文件是否存在
		f, err := staticFS.Open(path[1:]) // 去掉前导 /
		if err == nil {
			f.Close()
			staticServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// SPA路由回退到index.html
		c.Request.URL.Path = "/"
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
}
