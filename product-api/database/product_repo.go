package database

import (
	"fmt"
	"generated-api/models"

	"gorm.io/gorm"
)

// ProductRepository 商品数据访问层
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建仓库实例
func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: GetDB()}
}

// Create 创建商品
func (r *ProductRepository) Create(entity *models.Product) error {
	result := r.db.Create(entity)
	if result.Error != nil {
		return fmt.Errorf("创建商品失败: %w", result.Error)
	}
	return nil
}

// GetByID 根据ID查询商品
func (r *ProductRepository) GetByID(id int64) (*models.Product, error) {
	var entity models.Product
	result := r.db.First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询商品失败: %w", result.Error)
	}
	return &entity, nil
}

// List 分页查询商品列表
func (r *ProductRepository) List(params models.QueryProductParams) ([]models.Product, int64, error) {
	var entities []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// 关键字搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("sku LIKE ? OR name LIKE ? OR description LIKE ? OR image_url LIKE ?", keyword, keyword, keyword, keyword)
	}

	// 统计总数
	query.Count(&total)

	// 排序
	if params.OrderBy != "" {
		order := params.OrderBy
		if params.Order == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("id DESC")
	}

	// 分页
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	offset := (params.Page - 1) * params.PageSize
	result := query.Offset(offset).Limit(params.PageSize).Find(&entities)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("查询商品列表失败: %w", result.Error)
	}

	return entities, total, nil
}

// Update 更新商品
func (r *ProductRepository) Update(id int64, updates map[string]interface{}) error {
	result := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新商品失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("商品不存在")
	}
	return nil
}

// Delete 删除商品
func (r *ProductRepository) Delete(id int64) error {
	result := r.db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除商品失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("商品不存在")
	}
	return nil
}

// BatchDelete 批量删除商品
func (r *ProductRepository) BatchDelete(ids []int64) error {
	result := r.db.Delete(&models.Product{}, ids)
	if result.Error != nil {
		return fmt.Errorf("批量删除商品失败: %w", result.Error)
	}
	return nil
}
