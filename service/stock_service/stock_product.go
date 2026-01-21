package stock_service

import (
	"github.com/EDDYCJY/go-gin-example/models/stock"
)

// 公共产品字段结构体
type BaseStockProductForm struct {
	Unit                    string  `json:"unit" binding:"required,max=32"`
	Name                    string  `json:"name" binding:"required,max=255"`
	StockCustomizeProductID int     `json:"stock_customize_product_id" binding:"required"`
	SkuKey                  string  `json:"sku_key" binding:"max=255"`
	IsConsumable            int     `json:"is_consumable" binding:"gte=0,lte=1"`
	NoCodeNum               float64 `json:"no_code_num" binding:"gte=0"`
	CodeNum                 float64 `json:"code_num" binding:"gte=0"`
	TotalNum                float64 `json:"total_num" binding:"gte=0"`
	IsComponent             int     `json:"is_component" binding:"gte=0,lte=1"`
}

type AddStockProductForm struct {
	BaseStockProductForm
}

type EditStockProductForm struct {
	ID int `json:"id" binding:"required"`
	BaseStockProductForm
}

// 业务层 StockProduct 对象
type StockProduct struct {
	ID                      int
	Unit                    string
	Name                    string
	StockCustomizeProductID int
	SkuKey                  string
	IsConsumable            int
	NoCodeNum               float64
	CodeNum                 float64
	TotalNum                float64
	IsComponent             int
}

// 查询条件
type StockProductQuery struct {
	Name                    string
	StockCustomizeProductID int
	PageNum                 int
	PageSize                int
}

type StockProductDelete struct {
	ID int `json:"id" binding:"required"`
}

// ConvertAddFormToStockProduct 转换添加表单
func ConvertAddFormToStockProduct(form AddStockProductForm) StockProduct {
	return StockProduct{
		Unit:                    form.Unit,
		Name:                    form.Name,
		StockCustomizeProductID: form.StockCustomizeProductID,
		SkuKey:                  form.SkuKey,
		IsConsumable:            form.IsConsumable,
		NoCodeNum:               form.NoCodeNum,
		CodeNum:                 form.CodeNum,
		TotalNum:                form.TotalNum,
		IsComponent:             form.IsComponent,
	}
}

// ConvertEditFormToStockProduct 转换编辑表单
func ConvertEditFormToStockProduct(form EditStockProductForm) StockProduct {
	return StockProduct{
		ID:                      form.ID,
		Unit:                    form.Unit,
		Name:                    form.Name,
		StockCustomizeProductID: form.StockCustomizeProductID,
		SkuKey:                  form.SkuKey,
		IsConsumable:            form.IsConsumable,
		NoCodeNum:               form.NoCodeNum,
		CodeNum:                 form.CodeNum,
		TotalNum:                form.TotalNum,
		IsComponent:             form.IsComponent,
	}
}

// toMap 转换为数据库操作的 map
func (sp *StockProduct) toMap() map[string]interface{} {
	return map[string]interface{}{
		"unit":                       sp.Unit,
		"name":                       sp.Name,
		"stock_customize_product_id": sp.StockCustomizeProductID,
		"sku_key":                    sp.SkuKey,
		"is_consumable":              sp.IsConsumable,
		"no_code_num":                sp.NoCodeNum,
		"code_num":                   sp.CodeNum,
		"total_num":                  sp.TotalNum,
		"is_component":               sp.IsComponent,
	}
}

// Add 创建产品
func (sp *StockProduct) Add() error {
	return stock.AddStockProduct(sp.toMap())
}

// Edit 修改产品
func (sp *StockProduct) Edit() error {
	return stock.EditStockProduct(sp.ID, sp.toMap())
}

// GetStockProductByID 获取单个产品
func GetStockProductByID(id int) (*stock.StockProduct, error) {
	return stock.GetStockProduct(id)
}

// GetStockProductWithDetails 获取产品及其明细
func GetStockProductWithDetails(id int) (*stock.StockProduct, error) {
	return stock.GetStockProductWithDetails(id)
}

// GetAll 获取产品列表
func (q *StockProductQuery) GetAll() ([]*stock.StockProduct, error) {
	return stock.GetStockProducts(q.PageNum, q.PageSize, q.toMap())
}

func (q *StockProductQuery) GetAll1() ([]*stock.StockProduct, error) {
	return stock.GetStockProducts1(q.PageNum, q.PageSize, q.toMap())
}

// Count 获取产品总数
func (q *StockProductQuery) Count() (int, error) {
	return stock.GetStockProductTotal(q.toMap())
}

// toMap 查询条件转换
func (q *StockProductQuery) toMap() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if q.Name != "" {
		maps["name"] = q.Name
	}
	if q.StockCustomizeProductID > 0 {
		maps["stock_customize_product_id"] = q.StockCustomizeProductID
	}
	return maps
}

// Delete 删除产品
func (d *StockProductDelete) Delete() error {
	return stock.DeleteStockProduct(d.ID)
}

// ExistStockProductByID 检查产品是否存在
func ExistStockProductByID(id int) (bool, error) {
	return stock.ExistStockProductByID(id)
}

// ExistStockProductByName 检查产品名称是否存在
func ExistStockProductByName(name string) (bool, error) {
	return stock.ExistStockProductByName(name)
}

// UpdateStockProductNum 更新产品库存数量
func UpdateStockProductNum(id int, noCodeNum, codeNum float64) error {
	return stock.UpdateStockProductNum(id, noCodeNum, codeNum)
}
