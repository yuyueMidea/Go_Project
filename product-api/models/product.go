package models

import "time"

// Product 商品
type Product struct {
	// 主键ID
	ID int64 `json:"id" gorm:"primaryKey;column:id;type:integer;autoIncrement;not null;comment:主键ID"`
	// 商品编号
	Sku string `json:"sku" gorm:"column:sku;type:varchar(32);uniqueIndex;not null;comment:商品编号" binding:"required,max=32"`
	// 商品名称
	Name string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:商品名称" binding:"required,max=100"`
	// 商品描述
	Description string `json:"description" gorm:"column:description;type:text;comment:商品描述"`
	// 价格
	Price float64 `json:"price" gorm:"column:price;type:real;not null;comment:价格" binding:"required"`
	// 库存数量
	Stock int64 `json:"stock" gorm:"column:stock;type:integer;not null;comment:库存数量" binding:"required"`
	// 商品图片
	ImageURL string `json:"image_url" gorm:"column:image_url;type:varchar(500);comment:商品图片" binding:"url,max=500"`
	// 是否上架
	IsOnSale bool `json:"is_on_sale" gorm:"column:is_on_sale;type:boolean;not null;comment:是否上架" binding:"required"`
	// 重量(kg)
	Weight float64 `json:"weight" gorm:"column:weight;type:real;comment:重量(kg)"`
	// 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	// 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "product"
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Sku string `json:"sku" gorm:"column:sku;type:varchar(32);uniqueIndex;not null;comment:商品编号" binding:"required,max=32"`
	Name string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:商品名称" binding:"required,max=100"`
	Description string `json:"description" gorm:"column:description;type:text;comment:商品描述"`
	Price float64 `json:"price" gorm:"column:price;type:real;not null;comment:价格" binding:"required"`
	Stock int64 `json:"stock" gorm:"column:stock;type:integer;not null;comment:库存数量" binding:"required"`
	ImageURL string `json:"image_url" gorm:"column:image_url;type:varchar(500);comment:商品图片" binding:"url,max=500"`
	IsOnSale bool `json:"is_on_sale" gorm:"column:is_on_sale;type:boolean;not null;comment:是否上架" binding:"required"`
	Weight float64 `json:"weight" gorm:"column:weight;type:real;comment:重量(kg)"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	Sku string `json:"sku"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price *float64 `json:"price"`
	Stock *int64 `json:"stock"`
	ImageURL string `json:"image_url"`
	IsOnSale *bool `json:"is_on_sale"`
	Weight *float64 `json:"weight"`
}

// QueryProductParams 查询商品参数
type QueryProductParams struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	OrderBy  string `form:"order_by" json:"order_by"`
	Order    string `form:"order" json:"order"`
	Keyword  string `form:"keyword" json:"keyword"`
}

