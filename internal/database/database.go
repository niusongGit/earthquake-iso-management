package database

import (
	"earthquake-iso-management/internal/config"
	"earthquake-iso-management/internal/logger"
	"earthquake-iso-management/internal/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open(mysql.Open(config.C.MySQL.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.C.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.C.MySQL.MaxOpenConns)

	// 自动迁移
	if err := DB.AutoMigrate(&model.Admin{}, &model.Document{}); err != nil {
		return err
	}

	// 初始化默认管理员账号
	initAdmin()

	logger.Log.Info("数据库初始化完成")
	return nil
}

func initAdmin() {
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(config.C.Admin.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Log.Error("密码加密失败", zap.Error(err))
			return
		}
		admin := model.Admin{
			Username: config.C.Admin.Username,
			Password: string(hashedPassword),
		}
		if err := DB.Create(&admin).Error; err != nil {
			logger.Log.Error("创建默认管理员失败", zap.Error(err))
			return
		}
		logger.Log.Info("默认管理员账号创建成功", zap.String("username", config.C.Admin.Username))
	}
}
