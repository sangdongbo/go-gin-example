package order_service

import (
	"github.com/EDDYCJY/go-gin-example/models"
)

// 公共订单字段结构体
type BaseOrderForm struct {
	OrderSn         string  `json:"order_sn" binding:"required,max=64"`
	UserID          int     `json:"user_id" binding:"required"`
	ProductID       uint64  `json:"product_id" binding:"required"`
	ProductName     string  `json:"product_name" binding:"required,max=255"`
	ProductPrice    float64 `json:"product_price" binding:"required"`
	Quantity        uint    `json:"quantity" binding:"required,min=1"`
	TotalAmount     float64 `json:"total_amount" binding:"required"`
	DiscountAmount  float64 `json:"discount_amount" binding:"omitempty,min=0"`
	PayAmount       float64 `json:"pay_amount" binding:"required"`
	OrderStatus     uint8   `json:"order_status" binding:"gte=0,lte=4"`
	PaymentMethod   *uint8  `json:"payment_method" binding:"omitempty,gte=1,lte=3"`
	PaymentTime     string  `json:"payment_time"`
	ShippingAddress string  `json:"shipping_address" binding:"omitempty,max=512"`
	ShippingTime    string  `json:"shipping_time"`
	CompletionTime  string  `json:"completion_time"`
	CancelTime      string  `json:"cancel_time"`
}

type AddOrderForm struct {
	BaseOrderForm
}

type EditOrderForm struct {
	Id int `json:"id" binding:"required"`
	BaseOrderForm
}

// 业务层 Order 对象
type Order struct {
	Id              int
	OrderSn         string
	UserID          int
	ProductID       uint64
	ProductName     string
	ProductPrice    float64
	Quantity        uint
	TotalAmount     float64
	DiscountAmount  float64
	PayAmount       float64
	OrderStatus     uint8
	PaymentMethod   *uint8
	PaymentTime     string
	ShippingAddress string
	ShippingTime    string
	CompletionTime  string
	CancelTime      string
}

// 分页条件单独结构
type OrderQuery struct {
	OrderSn  string
	PageNum  int
	PageSize int
}

type OrderDelete struct {
	OrderSn string `json:"order_sn" binding:"required,max=64"`
}

func ConvertAddFormToOrder(form AddOrderForm) Order {
	order := Order{
		OrderSn:         form.OrderSn,
		UserID:          form.UserID,
		ProductID:       form.ProductID,
		ProductName:     form.ProductName,
		ProductPrice:    form.ProductPrice,
		Quantity:        form.Quantity,
		TotalAmount:     form.TotalAmount,
		DiscountAmount:  form.DiscountAmount,
		PayAmount:       form.PayAmount,
		OrderStatus:     form.OrderStatus,
		PaymentMethod:   form.PaymentMethod,
		PaymentTime:     form.PaymentTime,
		ShippingAddress: form.ShippingAddress,
		ShippingTime:    form.ShippingTime,
		CompletionTime:  form.CompletionTime,
		CancelTime:      form.CancelTime,
	}
	return order
}

// ConvertFormToOrder 将 Add/Edit 表单转换为 Order 实体
func ConvertEditFormToOrder(form EditOrderForm) Order {
	order := Order{
		Id:              form.Id,
		OrderSn:         form.OrderSn,
		UserID:          form.UserID,
		ProductID:       form.ProductID,
		ProductName:     form.ProductName,
		ProductPrice:    form.ProductPrice,
		Quantity:        form.Quantity,
		TotalAmount:     form.TotalAmount,
		DiscountAmount:  form.DiscountAmount,
		PayAmount:       form.PayAmount,
		OrderStatus:     form.OrderStatus,
		PaymentMethod:   form.PaymentMethod,
		PaymentTime:     form.PaymentTime,
		ShippingAddress: form.ShippingAddress,
		ShippingTime:    form.ShippingTime,
		CompletionTime:  form.CompletionTime,
		CancelTime:      form.CancelTime,
	}
	return order
}

// toModelOrder 转换为 models.Order（用于 DB 操作）
func toModelOrder(o *Order) *models.Order {
	model := &models.Order{
		Id:              o.Id,
		OrderSn:         o.OrderSn,
		UserId:          o.UserID,
		ProductId:       int(o.ProductID),
		ProductName:     o.ProductName,
		ProductPrice:    o.ProductPrice,
		Quantity:        int(o.Quantity),
		TotalAmount:     o.TotalAmount,
		DiscountAmount:  o.DiscountAmount,
		PayAmount:       o.PayAmount,
		OrderStatus:     int(o.OrderStatus),
		ShippingAddress: o.ShippingAddress,
	}
	if o.PaymentMethod != nil {
		model.PaymentMethod = int(*o.PaymentMethod)
	}
	return model
}

// Add 创建订单
func (o *Order) Add() error {
	return models.AddOrder(toModelOrder(o))
}

// Edit 修改订单
func (o *Order) Edit() error {
	return models.EditOrder(toModelOrder(o))
}

// GetOne 获取单个订单
func GetOne(id int) (models.Order, error) {
	return models.GetOrderById(id)
}

// GetAll 获取订单列表
func (q *OrderQuery) GetAll() ([]models.Order, error) {
	return models.GetOrdersWithProducts(q.PageNum, q.PageSize, q.toMap())
}

// Count 获取订单总数
func (q *OrderQuery) Count() (int, error) {
	return models.GetOrderTotal(q.toMap())
}

// 查询条件 map 构造器
func (q *OrderQuery) toMap() map[string]interface{} {
	maps := make(map[string]interface{})
	if q.OrderSn != "" {
		maps["order_sn"] = q.OrderSn
	}
	return maps
}

func (d *OrderDelete) Delete() error {
	return models.DeleteOrder(d.OrderSn)
}
