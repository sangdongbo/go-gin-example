package models

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	Model // 假设 Model 里有 ID、CreatedAt、UpdatedAt、DeletedAt

	ProductSn   string  `json:"product_sn" gorm:"type:varchar(50);not null;unique;comment:'商品编号'"`
	Name        string  `json:"name" gorm:"type:varchar(255);not null;comment:'商品名称'"`
	Description string  `json:"description" gorm:"type:text;comment:'商品描述'"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null;comment:'商品价格'"`
	Stock       int     `json:"stock" gorm:"default:0;comment:'商品库存数量'"`
}

// GetProducts 获取产品列表
func GetProducts(pageNum int, pageSize int, maps interface{}) ([]Product, error) {
	var (
		products []Product
		err      error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&products).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&products).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return products, nil
}
