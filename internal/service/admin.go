package service

import (
	"earthquake-iso-management/internal/database"
	"earthquake-iso-management/internal/logger"
	"earthquake-iso-management/internal/middleware"
	"earthquake-iso-management/internal/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务
type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// Login 管理员登录
func (s *AdminService) Login(req model.LoginRequest) (string, error) {
	var admin model.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		logger.Log.Warn("登录失败：用户不存在", zap.String("username", req.Username))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		logger.Log.Warn("登录失败：密码错误", zap.String("username", req.Username))
		return "", err
	}

	token, err := middleware.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		logger.Log.Error("生成令牌失败", zap.Error(err))
		return "", err
	}

	logger.Log.Info("管理员登录成功", zap.String("username", req.Username))
	return token, nil
}
