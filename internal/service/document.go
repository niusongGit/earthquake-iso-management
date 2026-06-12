package service

import (
	"earthquake-iso-management/internal/database"
	"earthquake-iso-management/internal/logger"
	"earthquake-iso-management/internal/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// DocumentService 文档服务
type DocumentService struct{}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

// GetDocumentList 获取文档列表（前台）
func (s *DocumentService) GetDocumentList(query model.DocumentListQuery) (*model.DocumentListResult, error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	db := database.DB.Model(&model.Document{})

	// 关键词搜索
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("iso_code LIKE ? OR name LIKE ? OR summary LIKE ? OR scope LIKE ? OR type LIKE ? OR standard_belongs_to LIKE ? OR belongs_to LIKE ? OR first_publish_code LIKE ?",
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// 阶段筛选
	if query.Stage != "" {
		db = db.Where("current_stage = ?", query.Stage)
	}

	// 排序
	switch query.Sort {
	case "publishDate":
		db = db.Order("publish_date DESC")
	default:
		db = db.Order("created_at DESC")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var documents []model.Document
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&documents).Error; err != nil {
		return nil, err
	}

	return &model.DocumentListResult{
		List:  documents,
		Total: total,
	}, nil
}

// GetDocumentByID 根据ID获取文档详情
func (s *DocumentService) GetDocumentByID(id uint) (*model.Document, error) {
	var doc model.Document
	if err := database.DB.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

// CreateDocument 创建文档
func (s *DocumentService) CreateDocument(req model.DocumentCreateRequest, file *multipart.FileHeader) (*model.Document, error) {
	doc := model.Document{
		ISOCode:             strings.TrimSpace(req.ISOCode),
		Name:                req.Name,
		Type:                req.Type,
		StandardBelongsTo:   req.StandardBelongsTo,
		BelongsTo:           req.BelongsTo,
		Summary:             req.Summary,
		Scope:               req.Scope,
		FirstPublishCode:    req.FirstPublishCode,
		CurrentStage:        req.CurrentStage,
		EarthquakeRelevance: req.EarthquakeRelevance,
	}

	// 解析发布日期
	if req.PublishDate != "" {
		t, err := time.Parse("2006-01-02", req.PublishDate)
		if err == nil {
			doc.PublishDate = &t
		}
	}

	// 处理附件上传
	if file != nil {
		attachmentPath, err := saveUploadFile(file, doc.ISOCode)
		if err != nil {
			return nil, fmt.Errorf("保存附件失败: %w", err)
		}
		doc.Attachment = attachmentPath
	}

	if err := database.DB.Create(&doc).Error; err != nil {
		return nil, err
	}

	logger.Log.Info("创建文档成功", zap.String("iso_code", doc.ISOCode), zap.Uint("id", doc.ID))
	return &doc, nil
}

// UpdateDocument 更新文档
func (s *DocumentService) UpdateDocument(id uint, req model.DocumentUpdateRequest, file *multipart.FileHeader) (*model.Document, error) {
	var doc model.Document
	if err := database.DB.First(&doc, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.ISOCode != "" {
		updates["iso_code"] = strings.TrimSpace(req.ISOCode)
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.StandardBelongsTo != "" {
		updates["standard_belongs_to"] = req.StandardBelongsTo
	}
	if req.BelongsTo != "" {
		updates["belongs_to"] = req.BelongsTo
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Scope != "" {
		updates["scope"] = req.Scope
	}
	if req.FirstPublishCode != "" {
		updates["first_publish_code"] = req.FirstPublishCode
	}
	if req.CurrentStage != "" {
		updates["current_stage"] = req.CurrentStage
	}
	if req.EarthquakeRelevance > 0 {
		updates["earthquake_relevance"] = req.EarthquakeRelevance
	}
	if req.PublishDate != "" {
		t, err := time.Parse("2006-01-02", req.PublishDate)
		if err == nil {
			updates["publish_date"] = &t
		}
	}

	// 处理附件上传
	if file != nil {
		isoCode := doc.ISOCode
		if v, ok := updates["iso_code"]; ok {
			isoCode = v.(string)
		}
		attachmentPath, err := saveUploadFile(file, isoCode)
		if err != nil {
			return nil, fmt.Errorf("保存附件失败: %w", err)
		}
		updates["attachment"] = attachmentPath

		// 删除旧附件
		if doc.Attachment != "" {
			oldPath := filepath.Join("uploads", doc.Attachment)
			os.Remove(oldPath)
		}
	}

	if err := database.DB.Model(&doc).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询返回完整数据
	database.DB.First(&doc, id)
	logger.Log.Info("更新文档成功", zap.Uint("id", id))
	return &doc, nil
}

// DeleteDocument 删除文档
func (s *DocumentService) DeleteDocument(id uint) error {
	var doc model.Document
	if err := database.DB.First(&doc, id).Error; err != nil {
		return err
	}

	// 删除附件文件
	if doc.Attachment != "" {
		oldPath := filepath.Join("uploads", doc.Attachment)
		os.Remove(oldPath)
	}

	if err := database.DB.Delete(&doc).Error; err != nil {
		return err
	}

	logger.Log.Info("删除文档成功", zap.Uint("id", id), zap.String("iso_code", doc.ISOCode))
	return nil
}

// saveUploadFile 保存上传文件
func saveUploadFile(file *multipart.FileHeader, isoCode string) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 生成文件名：ISO编号_时间戳.扩展名
	ext := path.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%d%s", isoCode, time.Now().UnixMilli(), ext)
	filePath := filepath.Join("uploads", fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fileName, nil
}
