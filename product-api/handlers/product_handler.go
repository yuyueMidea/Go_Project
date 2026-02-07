package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"generated-api/database"
	"generated-api/models"
)

// ProductHandler 商品HTTP处理器
type ProductHandler struct {
	repo *database.ProductRepository
}

// NewProductHandler 创建处理器实例
func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		repo: database.NewProductRepository(),
	}
}

// Create 创建商品
// @Summary 创建商品
// @Tags Product
func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	entity := models.Product{
		Sku: req.Sku,
		Name: req.Name,
		Description: req.Description,
		Price: req.Price,
		Stock: req.Stock,
		ImageURL: req.ImageURL,
		IsOnSale: req.IsOnSale,
		Weight: req.Weight,
	}

	if err := h.repo.Create(&entity); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, entity)
}

// GetByID 根据ID获取商品
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的ID")
		return
	}

	entity, err := h.repo.GetByID(id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}
	if entity == nil {
		NotFound(c, "商品不存在")
		return
	}

	Success(c, entity)
}

// List 获取商品列表
func (h *ProductHandler) List(c *gin.Context) {
	var params models.QueryProductParams
	if err := c.ShouldBindQuery(&params); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	entities, total, err := h.repo.List(params)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessPage(c, entities, total, params.Page, params.PageSize)
}

// Update 更新商品
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的ID")
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 构建更新字段 map
	updates := make(map[string]interface{})
	if req.Sku != "" {
		updates["sku"] = req.Sku
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.IsOnSale != nil {
		updates["is_on_sale"] = *req.IsOnSale
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}

	if len(updates) == 0 {
		BadRequest(c, "没有需要更新的字段")
		return
	}

	if err := h.repo.Update(id, updates); err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessMessage(c, "更新成功")
}

// Delete 删除商品
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的ID")
		return
	}

	if err := h.repo.Delete(id); err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessMessage(c, "删除成功")
}

// BatchDelete 批量删除商品
func (h *ProductHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.repo.BatchDelete(req.IDs); err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessMessage(c, "批量删除成功")
}
