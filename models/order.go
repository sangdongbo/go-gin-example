package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	Model                      // 假设 Model 里有 ID、CreatedAt、UpdatedAt、DeletedAt
	Id              int        `json:"id"`
	OrderSn         string     `json:"order_sn" gorm:"type:varchar(50);not null;unique;comment:'订单编号'"`
	UserId          int        `json:"user_id" gorm:"not null;comment:'用户ID'"`
	ProductId       int        `json:"product_id" gorm:"not null;comment:'商品ID'"`
	Product         Product    `json:"product" gorm:"foreignkey:ProductId;association_foreignkey:ID;"`
	ProductName     string     `json:"product_name" gorm:"type:varchar(255);comment:'商品名称'"`
	ProductPrice    float64    `json:"product_price" gorm:"type:decimal(10,2);comment:'商品单价'"`
	Quantity        int        `json:"quantity" gorm:"default:1;comment:'购买数量'"`
	TotalAmount     float64    `json:"total_amount" gorm:"type:decimal(10,2);comment:'订单总金额'"`
	DiscountAmount  float64    `json:"discount_amount" gorm:"type:decimal(10,2);default:0.00;comment:'优惠金额'"`
	PayAmount       float64    `json:"pay_amount" gorm:"type:decimal(10,2);comment:'实际支付金额'"`
	OrderStatus     int        `json:"order_status" gorm:"default:0;comment:'订单状态 0-待支付 1-待发货 2-已发货 3-已完成 4-已取消'"`
	PaymentMethod   int        `json:"payment_method" gorm:"comment:'支付方式 1-支付宝 2-微信 3-银行卡'"`
	PaymentTime     *time.Time `json:"payment_time" gorm:"comment:'支付时间'"`
	ShippingAddress string     `json:"shipping_address" gorm:"type:varchar(255);comment:'收货地址'"`
	ShippingTime    *time.Time `json:"shipping_time" gorm:"comment:'发货时间'"`
	CompletionTime  *time.Time `json:"completion_time" gorm:"comment:'完成时间'"`
	CancelTime      *time.Time `json:"cancel_time" gorm:"comment:'取消时间'"`
	Deleted         int        `json:"deleted" gorm:"default:0;comment:'是否删除 0-正常 1-删除'"`
}

// GetOrders 获取订单列表
func GetOrders(pageNum int, pageSize int, maps interface{}) ([]Order, error) {
	var (
		orders []Order
		err    error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&orders).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&orders).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

func GetOrderById(id int) (Order, error) {
	var order Order
	err := db.Where("id = ?", id).First(&order).Error
	return order, err
}

// GetOrdersWithProducts 获取订单和关联的产品信息
func GetOrdersWithProducts(pageNum int, pageSize int, maps interface{}) ([]Order, error) {
	var (
		orders []Order
		err    error
	)

	// 使用 Preload 方法加载 Product 关联数据
	if pageSize > 0 {
		err = db.Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Where("price > ?", 0) // 假设你想查询价格大于 0 的产品
		}).Where(maps).Find(&orders).Error
	} else {
		err = db.Preload("Product").Where(maps).Find(&orders).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}

// GetOrderTotal 获取订单数量
func GetOrderTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Order{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// AddOrder 添加订单
func AddOrder(order *Order) error {
	return db.Create(order).Error
}

func EditOrder(order *Order) error {
	return db.Model(&Order{}).Where("id = ?", order.Id).Updates(order).Error
}

func DeleteOrder(orderSn string) error {
	return db.Where("order_sn = ?", orderSn).Delete(&Order{}).Error
}
