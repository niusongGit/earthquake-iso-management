package middleware

import (
	"bytes"
	"earthquake-iso-management/internal/database"
	"earthquake-iso-management/internal/logger"
	"earthquake-iso-management/internal/model"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLog 操作日志中间件，记录后台管理操作
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 在 c.Next() 之前捕获请求体（body 只能读一次）
		requestBody := captureRequestBody(c)

		c.Next()

		// 只记录写操作（POST/PUT/DELETE）
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			return
		}

		// 获取管理员信息（由 AuthRequired 中间件或登录处理器设置）
		adminID := c.GetUint("admin_id")
		username, _ := c.Get("username")
		usernameStr, _ := username.(string)

		// 操作描述
		action := getActionName(method, c.Request.URL.Path)

		log := model.OperationLog{
			AdminID:     adminID,
			Username:    usernameStr,
			Method:      method,
			Path:        c.Request.URL.Path,
			Action:      action,
			TargetID:    c.Param("id"),
			RequestBody: requestBody,
			IP:          c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			Status:      c.Writer.Status(),
			CostTime:    time.Since(start).Milliseconds(),
		}

		// 异步写入数据库，避免影响接口响应速度
		go func() {
			if err := database.DB.Create(&log).Error; err != nil {
				logger.Log.Error("记录操作日志失败", zap.Error(err))
			}
		}()
	}
}

// captureRequestBody 捕获请求数据
// - JSON 请求：返回格式化后的 JSON 字符串
// - multipart/form-data：返回表单字段 + 文件名（不含文件内容）
// - 其他表单：返回表单字段
func captureRequestBody(c *gin.Context) string {
	contentType := c.Request.Header.Get("Content-Type")

	// multipart/form-data：解析表单，文件只存文件名
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(50 << 20); err != nil {
			return ""
		}
		data := map[string]interface{}{}

		// 普通表单字段
		if c.Request.PostForm != nil {
			for k, v := range c.Request.PostForm {
				if len(v) == 1 {
					data[k] = v[0]
				} else {
					data[k] = v
				}
			}
		}
		// 文件字段：只存文件名
		if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
			for field, files := range c.Request.MultipartForm.File {
				names := make([]string, 0, len(files))
				for _, f := range files {
					names = append(names, f.Filename)
				}
				if len(names) == 1 {
					data[field] = names[0]
				} else {
					data[field] = names
				}
			}
		}
		b, _ := json.Marshal(data)
		return string(b)
	}

	// application/json 或其他：读取 body
	if c.Request.Body != nil {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return ""
		}
		// 将 body 放回，供后续 handler 使用
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 尝试格式化 JSON
		if strings.HasPrefix(contentType, "application/json") && len(bodyBytes) > 0 {
			var buf bytes.Buffer
			if err := json.Compact(&buf, bodyBytes); err == nil {
				return buf.String()
			}
		}
		return string(bodyBytes)
	}

	return ""
}

// getActionName 根据请求方法和路径生成操作描述
func getActionName(method, path string) string {
	// 去除 /api/admin 前缀
	path = strings.TrimPrefix(path, "/api/admin")

	switch {
	case path == "/login" && method == "POST":
		return "管理员登录"
	case path == "/password" && method == "PUT":
		return "修改密码"
	case strings.HasPrefix(path, "/documents"):
		switch method {
		case "POST":
			return "创建文档"
		case "PUT":
			return "更新文档"
		case "DELETE":
			return "删除文档"
		}
	}
	return method + " " + path
}
