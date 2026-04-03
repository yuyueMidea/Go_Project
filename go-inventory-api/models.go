package models

import (
	"time"

	"gorm.io/gorm"
)

// Product 商品表
type Product struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;size:128" binding:"required"`
	Price     float64        `json:"price" gorm:"not null;check:price>=0"`
	Stock     int            `json:"stock" gorm:"not null;default:0;check:stock>=0"`
	Unit      string         `json:"unit" gorm:"size:16;default:'件'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}

// StockLogType 流水类型
type StockLogType string

const (
	LogTypeIn  StockLogType = "IN"  // 入库
	LogTypeOut StockLogType = "OUT" // 出库
)

// StockLog 库存流水表
type StockLog struct {
	ID        uint         `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID uint         `json:"product_id" gorm:"not null;index"`
	Product   Product      `json:"product" gorm:"foreignKey:ProductID"`
	Type      StockLogType `json:"type" gorm:"not null;size:4"` // IN / OUT
	Quantity  int          `json:"quantity" gorm:"not null;check:quantity>0"`
	Remark    string       `json:"remark" gorm:"size:256"`
	CreatedAt time.Time    `json:"created_at"`
}

// ========== 请求 / 响应 DTO ==========

type CreateProductReq struct {
	Name  string  `json:"name" binding:"required,min=1,max=128"`
	Price float64 `json:"price" binding:"required,gte=0"`
	Stock int     `json:"stock" binding:"gte=0"`
	Unit  string  `json:"unit"`
}

type UpdateProductReq struct {
	Name  *string  `json:"name" binding:"omitempty,min=1,max=128"`
	Price *float64 `json:"price" binding:"omitempty,gte=0"`
	Unit  *string  `json:"unit"`
}

type StockInReq struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	Remark    string `json:"remark"`
}

type StockOutReq struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
	Remark    string `json:"remark"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(data any) Response {
	return Response{Code: 0, Message: "ok", Data: data}
}

func Fail(code int, msg string) Response {
	return Response{Code: code, Message: msg}
}
