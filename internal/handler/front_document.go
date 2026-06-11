package handler

import (
	"earthquake-iso-management/internal/model"
	"earthquake-iso-management/internal/response"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetFrontDocumentList 前台文档列表
func GetFrontDocumentList(c *gin.Context) {
	query := model.DocumentListQuery{
		Keyword:  c.Query("keyword"),
		Stage:    c.Query("stage"),
		Sort:     c.Query("sort"),
		Page:     0,
		PageSize: 0,
	}
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		query.Page = p
	}
	if ps, err := strconv.Atoi(c.Query("pageSize")); err == nil && ps > 0 {
		query.PageSize = ps
	}

	result, err := docService.GetDocumentList(query)
	if err != nil {
		response.InternalError(c, "查询文档列表失败")
		return
	}

	response.SuccessPage(c, result.List, result.Total, query.Page, query.PageSize)
}

// GetFrontDocumentByID 前台文档详情
func GetFrontDocumentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	doc, err := docService.GetDocumentByID(uint(id))
	if err != nil {
		response.NotFound(c, "文档不存在")
		return
	}

	response.Success(c, doc)
}

// PreviewAttachment 预览PDF附件（inline方式，浏览器内嵌显示）
func PreviewAttachment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	doc, err := docService.GetDocumentByID(uint(id))
	if err != nil {
		response.NotFound(c, "文档不存在")
		return
	}

	if doc.Attachment == "" {
		response.NotFound(c, "该文档没有附件")
		return
	}

	filePath := filepath.Join("uploads", doc.Attachment)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "附件文件不存在")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline")
	c.File(filePath)
}

// DownloadAttachment 下载PDF附件（attachment方式，触发浏览器下载）
func DownloadAttachment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	doc, err := docService.GetDocumentByID(uint(id))
	if err != nil {
		response.NotFound(c, "文档不存在")
		return
	}

	if doc.Attachment == "" {
		response.NotFound(c, "该文档没有附件")
		return
	}

	filePath := filepath.Join("uploads", doc.Attachment)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "附件文件不存在")
		return
	}

	c.FileAttachment(filePath, doc.Attachment)
}
