package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"my-api/database"
	"my-api/models"
)

// TodoHandler 待办事项HTTP处理器
type TodoHandler struct {
	repo *database.TodoRepository
}

// NewTodoHandler 创建处理器实例
func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		repo: database.NewTodoRepository(),
	}
}

// Create 创建待办事项
// @Summary 创建待办事项
// @Tags Todo
func (h *TodoHandler) Create(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	entity := models.Todo{
		Title: req.Title,
		Done: req.Done,
		Priority: req.Priority,
	}

	if err := h.repo.Create(&entity); err != nil {
		InternalError(c, err.Error())
		return
	}

	Success(c, entity)
}

// GetByID 根据ID获取待办事项
func (h *TodoHandler) GetByID(c *gin.Context) {
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
		NotFound(c, "待办事项不存在")
		return
	}

	Success(c, entity)
}

// List 获取待办事项列表
func (h *TodoHandler) List(c *gin.Context) {
	var params models.QueryTodoParams
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

// Update 更新待办事项
func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "无效的ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 构建更新字段 map
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Done != nil {
		updates["done"] = *req.Done
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
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

// Delete 删除待办事项
func (h *TodoHandler) Delete(c *gin.Context) {
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

// BatchDelete 批量删除待办事项
func (h *TodoHandler) BatchDelete(c *gin.Context) {
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
