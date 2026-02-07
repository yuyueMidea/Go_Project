package database

import (
	"fmt"
	"my-api/models"

	"gorm.io/gorm"
)

// TodoRepository 待办事项数据访问层
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository 创建仓库实例
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{db: GetDB()}
}

// Create 创建待办事项
func (r *TodoRepository) Create(entity *models.Todo) error {
	result := r.db.Create(entity)
	if result.Error != nil {
		return fmt.Errorf("创建待办事项失败: %w", result.Error)
	}
	return nil
}

// GetByID 根据ID查询待办事项
func (r *TodoRepository) GetByID(id int64) (*models.Todo, error) {
	var entity models.Todo
	result := r.db.First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询待办事项失败: %w", result.Error)
	}
	return &entity, nil
}

// List 分页查询待办事项列表
func (r *TodoRepository) List(params models.QueryTodoParams) ([]models.Todo, int64, error) {
	var entities []models.Todo
	var total int64

	query := r.db.Model(&models.Todo{})

	// 关键字搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("title LIKE ?", keyword)
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
		return nil, 0, fmt.Errorf("查询待办事项列表失败: %w", result.Error)
	}

	return entities, total, nil
}

// Update 更新待办事项
func (r *TodoRepository) Update(id int64, updates map[string]interface{}) error {
	result := r.db.Model(&models.Todo{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新待办事项失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("待办事项不存在")
	}
	return nil
}

// Delete 删除待办事项
func (r *TodoRepository) Delete(id int64) error {
	result := r.db.Delete(&models.Todo{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除待办事项失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("待办事项不存在")
	}
	return nil
}

// BatchDelete 批量删除待办事项
func (r *TodoRepository) BatchDelete(ids []int64) error {
	result := r.db.Delete(&models.Todo{}, ids)
	if result.Error != nil {
		return fmt.Errorf("批量删除待办事项失败: %w", result.Error)
	}
	return nil
}
