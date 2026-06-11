package model

import (
	"time"
)

// Document 国际标准文档模型
type Document struct {
	ID                  uint       `json:"id" gorm:"primaryKey;comment:主键ID"`
	ISOCode             string     `json:"iso_code" gorm:"size:50;not null;column:iso_code;comment:国际标准编号"`
	Name                string     `json:"name" gorm:"size:255;not null;comment:国际标准名称"`
	Type                string     `json:"type" gorm:"size:20;not null;comment:标准类型"`
	StandardBelongsTo   string     `json:"standard_belongs_to" gorm:"size:100;column:standard_belongs_to;comment:标准所属"`
	BelongsTo           string     `json:"belongs_to" gorm:"size:50;column:belongs_to;comment:所属"`
	Summary             string     `json:"summary" gorm:"type:text;comment:摘要"`
	Scope               string     `json:"scope" gorm:"type:text;comment:范围"`
	PublishDate         *time.Time `json:"publish_date" gorm:"type:date;comment:发布日期"`
	FirstPublishCode    string     `json:"first_publish_code" gorm:"size:100;comment:首次发布编号"`
	CurrentStage        string     `json:"current_stage" gorm:"size:20;not null;comment:当前阶段"`
	EarthquakeRelevance int        `json:"earthquake_relevance" gorm:"default:1;comment:地震相关度"`
	CreatedAt           time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	Attachment          string     `json:"attachment" gorm:"size:500;comment:附件路径"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "documents"
}

// DocumentListQuery 前台文档列表查询参数
type DocumentListQuery struct {
	Keyword  string `form:"keyword"`
	Stage    string `form:"stage"`
	Sort     string `form:"sort"` // createTime 或 publishDate
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
}

// DocumentListResult 文档列表返回结果
type DocumentListResult struct {
	List  []Document `json:"list"`
	Total int64      `json:"total"`
}

// DocumentCreateRequest 创建文档请求
type DocumentCreateRequest struct {
	ISOCode             string `json:"iso_code" form:"iso_code" binding:"required"`
	Name                string `json:"name" form:"name" binding:"required"`
	Type                string `json:"type" form:"type" binding:"required"`
	StandardBelongsTo   string `json:"standard_belongs_to" form:"standard_belongs_to"`
	BelongsTo           string `json:"belongs_to" form:"belongs_to"`
	Summary             string `json:"summary" form:"summary"`
	Scope               string `json:"scope" form:"scope"`
	PublishDate         string `json:"publish_date" form:"publish_date"`
	FirstPublishCode    string `json:"first_publish_code" form:"first_publish_code"`
	CurrentStage        string `json:"current_stage" form:"current_stage" binding:"required"`
	EarthquakeRelevance int    `json:"earthquake_relevance" form:"earthquake_relevance"`
}

// DocumentUpdateRequest 更新文档请求
type DocumentUpdateRequest struct {
	ISOCode             string `json:"iso_code" form:"iso_code"`
	Name                string `json:"name" form:"name"`
	Type                string `json:"type" form:"type"`
	StandardBelongsTo   string `json:"standard_belongs_to" form:"standard_belongs_to"`
	BelongsTo           string `json:"belongs_to" form:"belongs_to"`
	Summary             string `json:"summary" form:"summary"`
	Scope               string `json:"scope" form:"scope"`
	PublishDate         string `json:"publish_date" form:"publish_date"`
	FirstPublishCode    string `json:"first_publish_code" form:"first_publish_code"`
	CurrentStage        string `json:"current_stage" form:"current_stage"`
	EarthquakeRelevance int    `json:"earthquake_relevance" form:"earthquake_relevance"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}
