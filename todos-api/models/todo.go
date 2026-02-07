package models

import "time"

// Todo 待办事项
type Todo struct {
	// 主键ID
	ID int64 `json:"id" gorm:"primaryKey;column:id;type:integer;autoIncrement;not null;comment:主键ID"`
	// 标题
	Title string `json:"title" gorm:"column:title;type:varchar(200);not null;comment:标题" binding:"required,max=200"`
	// 是否完成
	Done bool `json:"done" gorm:"column:done;type:boolean;not null;comment:是否完成" binding:"required"`
	// 优先级: 0低 1中 2高
	Priority int64 `json:"priority" gorm:"column:priority;type:integer;comment:优先级: 0低 1中 2高"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Todo) TableName() string {
	return "todo"
}

// CreateTodoRequest 创建待办事项请求
type CreateTodoRequest struct {
	Title string `json:"title" gorm:"column:title;type:varchar(200);not null;comment:标题" binding:"required,max=200"`
	Done bool `json:"done" gorm:"column:done;type:boolean;not null;comment:是否完成" binding:"required"`
	Priority int64 `json:"priority" gorm:"column:priority;type:integer;comment:优先级: 0低 1中 2高"`
}

// UpdateTodoRequest 更新待办事项请求
type UpdateTodoRequest struct {
	Title string `json:"title"`
	Done *bool `json:"done"`
	Priority *int64 `json:"priority"`
}

// QueryTodoParams 查询待办事项参数
type QueryTodoParams struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	OrderBy  string `form:"order_by" json:"order_by"`
	Order    string `form:"order" json:"order"`
	Keyword  string `form:"keyword" json:"keyword"`
}

