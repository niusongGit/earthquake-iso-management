package handler

import (
	"earthquake-iso-management/internal/model"
	"earthquake-iso-management/internal/response"
	"earthquake-iso-management/internal/service"

	"github.com/gin-gonic/gin"
)

var adminService = service.NewAdminService()

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	token, err := adminService.Login(req)
	if err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	response.Success(c, model.LoginResponse{Token: token})
}
