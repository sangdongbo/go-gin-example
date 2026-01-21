package stock

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/jinzhu/gorm"
)

type StockProductDetail struct {
	models.Model

	StockProductID int     `json:"stock_product_id" gorm:"not null;default:0;index:idx_stock_product_id,idx_company_product" comment:"仓库产品id"`
	NeedReturn     int     `json:"need_return" gorm:"type:tinyint;not null;default:0" comment:"是否需归还 0=否 1=是"`
	Num            float64 `json:"num" gorm:"type:decimal(11,2);not null;default:0.00" comment:"数量"`
	OrderID        int     `json:"order_id" gorm:"not null;default:0;index:idx_order_id" comment:"订单id"`
	CostPrice      float64 `json:"cost_price" gorm:"type:decimal(10,2);not null;default:0.00" comment:"成本价"`
	Code           string  `json:"code" gorm:"type:varchar(32);not null;default:'';index:idx_code" comment:"序列号"`
	HiddenCode     string  `json:"hidden_code" gorm:"type:varchar(64);not null;default:'';unique_index:uk_spd_company_batch_hidden_code;index:idx_hidden_code" comment:"唯一编码（系统使用）"`
	Status         int     `json:"status" gorm:"type:tinyint;not null;default:0;index:idx_status" comment:"状态 0=正常 1=丢失 2=报废"`
	Note           string  `json:"note" gorm:"type:varchar(255);not null;default:''"`

	// 关联关系 - 属于某个产品
	StockProduct StockProduct `json:"stock_product" gorm:"foreignkey:StockProductID"`
}

// ExistStockProductDetailByID 根据ID检查明细是否存在
func ExistStockProductDetailByID(id int) (bool, error) {
	var detail StockProductDetail
	err := models.Db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&detail).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if detail.ID > 0 {
		return true, nil
	}

	return false, nil
}

// ExistStockProductDetailByHiddenCode 根据唯一编码检查明细是否存在
func ExistStockProductDetailByHiddenCode(hiddenCode string) (bool, error) {
	var detail StockProductDetail
	err := models.Db.Select("id").Where("hidden_code = ? AND deleted_on = ?", hiddenCode, 0).First(&detail).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if detail.ID > 0 {
		return true, nil
	}

	return false, nil
}

// GetStockProductDetailTotal 获取明细总数
func GetStockProductDetailTotal(maps interface{}) (int, error) {
	var count int
	if err := models.Db.Model(&StockProductDetail{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetStockProductDetails 获取明细列表（分页）
func GetStockProductDetails(pageNum int, pageSize int, maps interface{}) ([]*StockProductDetail, error) {
	var details []*StockProductDetail
	err := models.Db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&details).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return details, nil
}

// GetStockProductDetail 根据ID获取单个明细
func GetStockProductDetail(id int) (*StockProductDetail, error) {
	var detail StockProductDetail
	err := models.Db.Where("id = ? AND deleted_on = ?", id, 0).First(&detail).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &detail, nil
}

// GetStockProductDetailWithProduct 根据ID获取明细及关联的产品信息
func GetStockProductDetailWithProduct(id int) (*StockProductDetail, error) {
	var detail StockProductDetail
	err := models.Db.Where("id = ? AND deleted_on = ?", id, 0).First(&detail).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 预加载关联的产品数据
	err = models.Db.Model(&detail).Related(&detail.StockProduct).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &detail, nil
}

// GetStockProductDetailsByProductID 根据产品ID获取明细列表
func GetStockProductDetailsByProductID(productID int, pageNum int, pageSize int) ([]*StockProductDetail, error) {
	var details []*StockProductDetail
	err := models.Db.Where("stock_product_id = ? AND deleted_on = ?", productID, 0).
		Offset(pageNum).
		Limit(pageSize).
		Find(&details).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return details, nil
}

// GetStockProductDetailsByOrderID 根据订单ID获取明细列表
func GetStockProductDetailsByOrderID(orderID int) ([]*StockProductDetail, error) {
	var details []*StockProductDetail
	err := models.Db.Where("order_id = ? AND deleted_on = ?", orderID, 0).Find(&details).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return details, nil
}

// GetStockProductDetailsByStatus 根据状态获取明细列表
func GetStockProductDetailsByStatus(status int, pageNum int, pageSize int) ([]*StockProductDetail, error) {
	var details []*StockProductDetail
	err := models.Db.Where("status = ? AND deleted_on = ?", status, 0).
		Offset(pageNum).
		Limit(pageSize).
		Find(&details).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return details, nil
}

// AddStockProductDetail 添加明细
func AddStockProductDetail(data map[string]interface{}) error {
	detail := StockProductDetail{
		StockProductID: data["stock_product_id"].(int),
		NeedReturn:     data["need_return"].(int),
		Num:            data["num"].(float64),
		OrderID:        data["order_id"].(int),
		CostPrice:      data["cost_price"].(float64),
		Code:           data["code"].(string),
		HiddenCode:     data["hidden_code"].(string),
		Status:         data["status"].(int),
		Note:           data["note"].(string),
	}

	if err := models.Db.Create(&detail).Error; err != nil {
		return err
	}

	return nil
}

// EditStockProductDetail 修改明细
func EditStockProductDetail(id int, data interface{}) error {
	if err := models.Db.Model(&StockProductDetail{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteStockProductDetail 删除明细（软删除）
func DeleteStockProductDetail(id int) error {
	if err := models.Db.Where("id = ?", id).Delete(StockProductDetail{}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateStockProductDetailStatus 更新明细状态
func UpdateStockProductDetailStatus(id int, status int, note string) error {
	updates := map[string]interface{}{
		"status": status,
		"note":   note,
	}

	if err := models.Db.Model(&StockProductDetail{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// BatchAddStockProductDetails 批量添加明细
func BatchAddStockProductDetails(details []StockProductDetail) error {
	tx := models.Db.Begin()

	for _, detail := range details {
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetStockProductDetailSummary 获取明细汇总信息（按产品ID）
func GetStockProductDetailSummary(productID int) (map[string]interface{}, error) {
	type Summary struct {
		TotalNum      float64
		NormalNum     float64
		LostNum       float64
		ScrapNum      float64
		TotalValue    float64
		NeedReturnNum float64
	}

	var summary Summary
	err := models.Db.Model(&StockProductDetail{}).
		Select("SUM(num) as total_num, "+
			"SUM(CASE WHEN status = 0 THEN num ELSE 0 END) as normal_num, "+
			"SUM(CASE WHEN status = 1 THEN num ELSE 0 END) as lost_num, "+
			"SUM(CASE WHEN status = 2 THEN num ELSE 0 END) as scrap_num, "+
			"SUM(num * cost_price) as total_value, "+
			"SUM(CASE WHEN need_return = 1 THEN num ELSE 0 END) as need_return_num").
		Where("stock_product_id = ? AND deleted_on = ?", productID, 0).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total_num":       summary.TotalNum,
		"normal_num":      summary.NormalNum,
		"lost_num":        summary.LostNum,
		"scrap_num":       summary.ScrapNum,
		"total_value":     summary.TotalValue,
		"need_return_num": summary.NeedReturnNum,
	}

	return result, nil
}
