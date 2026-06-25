package service

import (
	"earthquake-iso-management/internal/database"
	"earthquake-iso-management/internal/logger"
	"earthquake-iso-management/internal/middleware"
	"earthquake-iso-management/internal/model"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务
type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// Login 管理员登录
func (s *AdminService) Login(req model.LoginRequest) (string, model.Admin, error) {
	var admin model.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		logger.Log.Warn("登录失败：用户不存在", zap.String("username", req.Username))
		return "", model.Admin{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		logger.Log.Warn("登录失败：密码错误", zap.String("username", req.Username))
		return "", model.Admin{}, err
	}

	token, err := middleware.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		logger.Log.Error("生成令牌失败", zap.Error(err))
		return "", model.Admin{}, err
	}

	logger.Log.Info("管理员登录成功", zap.String("username", req.Username))
	return token, admin, nil
}

// ChangePassword 修改密码
func (s *AdminService) ChangePassword(adminID uint, req model.ChangePasswordRequest) error {
	var admin model.Admin
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		logger.Log.Warn("修改密码失败：用户不存在", zap.Uint("admin_id", adminID))
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		logger.Log.Warn("修改密码失败：旧密码错误", zap.Uint("admin_id", adminID))
		return fmt.Errorf("旧密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error("密码加密失败", zap.Error(err))
		return err
	}

	if err := database.DB.Model(&admin).Update("password", string(hashedPassword)).Error; err != nil {
		logger.Log.Error("更新密码失败", zap.Error(err))
		return err
	}

	logger.Log.Info("管理员修改密码成功", zap.Uint("admin_id", adminID))
	return nil
}
