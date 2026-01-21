package stock_service

import (
	"github.com/EDDYCJY/go-gin-example/models/stock"
)

// 公共明细字段结构体
type BaseStockProductDetailForm struct {
	StockProductID int     `json:"stock_product_id" binding:"required"`
	NeedReturn     int     `json:"need_return" binding:"gte=0,lte=1"`
	Num            float64 `json:"num" binding:"required,gt=0"`
	OrderID        int     `json:"order_id"`
	CostPrice      float64 `json:"cost_price" binding:"gte=0"`
	Code           string  `json:"code" binding:"max=32"`
	HiddenCode     string  `json:"hidden_code" binding:"required,max=64"`
	Status         int     `json:"status" binding:"gte=0,lte=2"`
	Note           string  `json:"note" binding:"max=255"`
}

type AddStockProductDetailForm struct {
	BaseStockProductDetailForm
}

type EditStockProductDetailForm struct {
	ID int `json:"id" binding:"required"`
	BaseStockProductDetailForm
}

// 业务层 StockProductDetail 对象
type StockProductDetail struct {
	ID             int
	StockProductID int
	NeedReturn     int
	Num            float64
	OrderID        int
	CostPrice      float64
	Code           string
	HiddenCode     string
	Status         int
	Note           string
}

// 查询条件
type StockProductDetailQuery struct {
	StockProductID int
	OrderID        int
	Status         int
	PageNum        int
	PageSize       int
}

type StockProductDetailDelete struct {
	ID int `json:"id" binding:"required"`
}

// ConvertAddFormToStockProductDetail 转换添加表单
func ConvertAddFormToStockProductDetail(form AddStockProductDetailForm) StockProductDetail {
	return StockProductDetail{
		StockProductID: form.StockProductID,
		NeedReturn:     form.NeedReturn,
		Num:            form.Num,
		OrderID:        form.OrderID,
		CostPrice:      form.CostPrice,
		Code:           form.Code,
		HiddenCode:     form.HiddenCode,
		Status:         form.Status,
		Note:           form.Note,
	}
}

// ConvertEditFormToStockProductDetail 转换编辑表单
func ConvertEditFormToStockProductDetail(form EditStockProductDetailForm) StockProductDetail {
	return StockProductDetail{
		ID:             form.ID,
		StockProductID: form.StockProductID,
		NeedReturn:     form.NeedReturn,
		Num:            form.Num,
		OrderID:        form.OrderID,
		CostPrice:      form.CostPrice,
		Code:           form.Code,
		HiddenCode:     form.HiddenCode,
		Status:         form.Status,
		Note:           form.Note,
	}
}

// toMap 转换为数据库操作的 map
func (spd *StockProductDetail) toMap() map[string]interface{} {
	return map[string]interface{}{
		"stock_product_id": spd.StockProductID,
		"need_return":      spd.NeedReturn,
		"num":              spd.Num,
		"order_id":         spd.OrderID,
		"cost_price":       spd.CostPrice,
		"code":             spd.Code,
		"hidden_code":      spd.HiddenCode,
		"status":           spd.Status,
		"note":             spd.Note,
	}
}

// Add 创建明细
func (spd *StockProductDetail) Add() error {
	return stock.AddStockProductDetail(spd.toMap())
}

// Edit 修改明细
func (spd *StockProductDetail) Edit() error {
	return stock.EditStockProductDetail(spd.ID, spd.toMap())
}

// GetStockProductDetailByID 获取单个明细
func GetStockProductDetailByID(id int) (*stock.StockProductDetail, error) {
	return stock.GetStockProductDetail(id)
}

// GetStockProductDetailWithProduct 获取明细及关联产品
func GetStockProductDetailWithProduct(id int) (*stock.StockProductDetail, error) {
	return stock.GetStockProductDetailWithProduct(id)
}

// GetAll 获取明细列表
func (q *StockProductDetailQuery) GetAll() ([]*stock.StockProductDetail, error) {
	return stock.GetStockProductDetails(q.PageNum, q.PageSize, q.toMap())
}

// Count 获取明细总数
func (q *StockProductDetailQuery) Count() (int, error) {
	return stock.GetStockProductDetailTotal(q.toMap())
}

// toMap 查询条件转换
func (q *StockProductDetailQuery) toMap() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if q.StockProductID > 0 {
		maps["stock_product_id"] = q.StockProductID
	}
	if q.OrderID > 0 {
		maps["order_id"] = q.OrderID
	}
	if q.Status >= 0 {
		maps["status"] = q.Status
	}
	return maps
}

// Delete 删除明细
func (d *StockProductDetailDelete) Delete() error {
	return stock.DeleteStockProductDetail(d.ID)
}

// ExistStockProductDetailByID 检查明细是否存在
func ExistStockProductDetailByID(id int) (bool, error) {
	return stock.ExistStockProductDetailByID(id)
}

// ExistStockProductDetailByHiddenCode 检查唯一编码是否存在
func ExistStockProductDetailByHiddenCode(hiddenCode string) (bool, error) {
	return stock.ExistStockProductDetailByHiddenCode(hiddenCode)
}

// UpdateStockProductDetailStatus 更新明细状态
func UpdateStockProductDetailStatus(id int, status int, note string) error {
	return stock.UpdateStockProductDetailStatus(id, status, note)
}

// GetStockProductDetailSummary 获取明细汇总信息
func GetStockProductDetailSummary(productID int) (map[string]interface{}, error) {
	return stock.GetStockProductDetailSummary(productID)
}
