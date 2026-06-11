package handler

import (
	"earthquake-iso-management/internal/model"
	"earthquake-iso-management/internal/response"
	"earthquake-iso-management/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var docService = service.NewDocumentService()

// CreateDocument 创建文档
func CreateDocument(c *gin.Context) {
	var req model.DocumentCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// StandardBelongsTo默认值
	if req.StandardBelongsTo == "" {
		req.StandardBelongsTo = "ISO/TC"
	}
	// 地震相关度默认值
	if req.EarthquakeRelevance <= 0 {
		req.EarthquakeRelevance = 1
	}

	file, _ := c.FormFile("attachment")
	doc, err := docService.CreateDocument(req, file)
	if err != nil {
		response.InternalError(c, "创建文档失败: "+err.Error())
		return
	}

	response.Success(c, doc)
}

// GetAdminDocumentList 后台文档列表
func GetAdminDocumentList(c *gin.Context) {
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

// GetAdminDocumentByID 后台获取文档详情
func GetAdminDocumentByID(c *gin.Context) {
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

// UpdateDocument 更新文档
func UpdateDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	var req model.DocumentUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	file, _ := c.FormFile("attachment")
	doc, err := docService.UpdateDocument(uint(id), req, file)
	if err != nil {
		response.InternalError(c, "更新文档失败: "+err.Error())
		return
	}

	response.Success(c, doc)
}

// DeleteDocument 删除文档
func DeleteDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的文档ID")
		return
	}

	if err := docService.DeleteDocument(uint(id)); err != nil {
		response.NotFound(c, "文档不存在")
		return
	}

	response.Success(c, nil)
}
