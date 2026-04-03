package handlers

import (
	"net/http"
	"strconv"

	"inventory-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// CreateProduct POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(400, err.Error()))
		return
	}

	product := models.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
		Unit:  req.Unit,
	}
	if product.Unit == "" {
		product.Unit = "件"
	}

	if err := h.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusConflict, models.Fail(409, "商品名称已存在或数据异常: "+err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.OK(product))
}

// ListProducts GET /products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var products []models.Product
	query := h.db.Model(&models.Product{})

	// 支持按名称模糊搜索 ?name=xxx
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	// 支持分页 ?page=1&size=20
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var total int64
	query.Count(&total)
	query.Offset((page - 1) * size).Limit(size).Find(&products)

	c.JSON(http.StatusOK, models.OK(gin.H{
		"total":    total,
		"page":     page,
		"size":     size,
		"products": products,
	}))
}

// GetProduct GET /products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	product, ok := h.findProduct(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, models.OK(product))
}

// UpdateProduct PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	product, ok := h.findProduct(c)
	if !ok {
		return
	}

	var req models.UpdateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(400, err.Error()))
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Unit != nil {
		updates["unit"] = *req.Unit
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, models.Fail(400, "没有需要更新的字段"))
		return
	}

	if err := h.db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.OK(product))
}

// DeleteProduct DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	product, ok := h.findProduct(c)
	if !ok {
		return
	}

	// 软删除
	if err := h.db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Fail(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.OK(nil))
}

// ─── helper ────────────────────────────────────────────────────────────────

func (h *ProductHandler) findProduct(c *gin.Context) (models.Product, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, models.Fail(400, "无效的商品 ID"))
		return models.Product{}, false
	}

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Fail(404, "商品不存在"))
		return models.Product{}, false
	}
	return product, true
}
