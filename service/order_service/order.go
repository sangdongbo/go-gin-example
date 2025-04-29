package order_service

import (
	"github.com/EDDYCJY/go-gin-example/models"
)

type AddOrderForm struct {
	OrderSn         string  `json:"order_sn" binding:"required,max=64"`             // 订单编号
	UserID          int     `json:"user_id" binding:"required"`                     // 下单用户ID
	ProductID       uint64  `json:"product_id" binding:"required"`                  // 商品ID
	ProductName     string  `json:"product_name" binding:"required,max=255"`        // 商品名称快照
	ProductPrice    float64 `json:"product_price" binding:"required"`               // 商品单价快照
	Quantity        uint    `json:"quantity" binding:"required,min=1"`              // 购买数量，最少1件
	TotalAmount     float64 `json:"total_amount" binding:"required"`                // 订单总金额
	DiscountAmount  float64 `json:"discount_amount" binding:"omitempty,min=0"`      // 优惠金额，允许空，最小0
	PayAmount       float64 `json:"pay_amount" binding:"required"`                  // 实际支付金额
	OrderStatus     uint8   `json:"order_status" binding:"gte=0,lte=4"`             // 订单状态 0-4
	PaymentMethod   *uint8  `json:"payment_method" binding:"omitempty,gte=1,lte=3"` // 支付方式，可空
	PaymentTime     string  `json:"payment_time"`                                   // 支付时间，可空
	ShippingAddress string  `json:"shipping_address" binding:"omitempty,max=512"`   // 收货地址，可空
	ShippingTime    string  `json:"shipping_time"`                                  // 发货时间，可空
	CompletionTime  string  `json:"completion_time"`                                // 完成时间，可空
	CancelTime      string  `json:"cancel_time"`                                    // 取消时间，可空
}

type Order struct {
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

	PageNum  int
	PageSize int
}

// ConvertAddOrderFormToOrder 转换 AddOrderForm 到 Order
func ConvertAddOrderFormToOrder(form AddOrderForm) Order {
	return Order{
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
}

// GetAll 获取所有订单
func (t *Order) GetAll() ([]models.Order, error) {
	var orders = []models.Order{}

	orders, err := models.GetOrdersWithProducts(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// Count 获取订单总数
func (t *Order) Count() (int, error) {
	return models.GetOrderTotal(t.getMaps())
}

func (t *Order) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if t.OrderSn != "" {
		maps["order_sn"] = t.OrderSn
	}

	return maps
}

// Add 添加订单
func (o *Order) Add() error {
	modelOrder := models.Order{
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
		PaymentMethod:   0, // 注意处理空指针 o.PaymentMethod
		ShippingAddress: o.ShippingAddress,
	}

	if o.PaymentMethod != nil {
		modelOrder.PaymentMethod = int(*o.PaymentMethod)
	}

	return models.AddOrder(&modelOrder)
}
