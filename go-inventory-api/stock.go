package handlers

import (
	"net/http"
	"strconv"

	"inventory-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockHandler struct {
	db *gorm.DB
}

func NewStockHandler(db *gorm.DB) *StockHandler {
	return &StockHandler{db: db}
}

// StockIn POST /stock/in  —— 入库
func (h *StockHandler) StockIn(c *gin.Context) {
	var req models.StockInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(400, err.Error()))
		return
	}

	var log models.StockLog
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var product models.Product
		// 加行锁（SQLite 使用事务隔离）
		if err := tx.First(&product, req.ProductID).Error; err != nil {
			return err
		}

		product.Stock += req.Quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		log = models.StockLog{
			ProductID: req.ProductID,
			Type:      models.LogTypeIn,
			Quantity:  req.Quantity,
			Remark:    req.Remark,
		}
		return tx.Create(&log).Error
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Fail(404, "商品不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, models.Fail(500, err.Error()))
		}
		return
	}

	// 预加载商品信息再返回
	h.db.Preload("Product").First(&log, log.ID)
	c.JSON(http.StatusOK, models.OK(log))
}

// StockOut POST /stock/out  —— 出库
func (h *StockHandler) StockOut(c *gin.Context) {
	var req models.StockOutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail(400, err.Error()))
		return
	}

	var log models.StockLog
	err := h.db.Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.First(&product, req.ProductID).Error; err != nil {
			return err
		}

		if product.Stock < req.Quantity {
			return &InsufficientStockError{Available: product.Stock, Requested: req.Quantity}
		}

		product.Stock -= req.Quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		log = models.StockLog{
			ProductID: req.ProductID,
			Type:      models.LogTypeOut,
			Quantity:  req.Quantity,
			Remark:    req.Remark,
		}
		return tx.Create(&log).Error
	})

	if err != nil {
		switch e := err.(type) {
		case *InsufficientStockError:
			c.JSON(http.StatusUnprocessableEntity, models.Fail(422,
				"库存不足，当前库存: "+strconv.Itoa(e.Available)+
					"，申请出库: "+strconv.Itoa(e.Requested)))
		default:
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, models.Fail(404, "商品不存在"))
			} else {
				c.JSON(http.StatusInternalServerError, models.Fail(500, err.Error()))
			}
		}
		return
	}

	h.db.Preload("Product").First(&log, log.ID)
	c.JSON(http.StatusOK, models.OK(log))
}

// ListStockLogs GET /stock/logs  —— 流水查询
func (h *StockHandler) ListStockLogs(c *gin.Context) {
	query := h.db.Model(&models.StockLog{}).Preload("Product")

	// 按商品过滤
	if pid := c.Query("product_id"); pid != "" {
		query = query.Where("product_id = ?", pid)
	}
	// 按类型过滤 ?type=IN 或 ?type=OUT
	if t := c.Query("type"); t == "IN" || t == "OUT" {
		query = query.Where("type = ?", t)
	}

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

	var logs []models.StockLog
	query.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&logs)

	c.JSON(http.StatusOK, models.OK(gin.H{
		"total": total,
		"page":  page,
		"size":  size,
		"logs":  logs,
	}))
}

// GetProductStock GET /stock/:product_id  —— 单商品库存
func (h *StockHandler) GetProductStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("product_id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, models.Fail(400, "无效的商品 ID"))
		return
	}

	var product models.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Fail(404, "商品不存在"))
		return
	}
	c.JSON(http.StatusOK, models.OK(gin.H{
		"product_id": product.ID,
		"name":       product.Name,
		"unit":       product.Unit,
		"stock":      product.Stock,
	}))
}

// ─── custom errors ──────────────────────────────────────────────────────────

type InsufficientStockError struct {
	Available int
	Requested int
}

func (e *InsufficientStockError) Error() string {
	return "库存不足"
}
