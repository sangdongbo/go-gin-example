package stock

import (
	"fmt"
	"time"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/jinzhu/gorm"
)

type StockProduct struct {
	models.Model

	Unit                    string  `json:"unit" gorm:"type:varchar(32);not null;default:''" comment:"单位"`
	Name                    string  `json:"name" gorm:"type:varchar(255);not null;default:'';index:idx_name" comment:"sku 名称（唯一）"`
	StockCustomizeProductID int     `json:"stock_customize_product_id" gorm:"not null;default:0;index:idx_stock_customize_product_id" comment:"公司产品ID"`
	SkuKey                  string  `json:"sku_key" gorm:"type:varchar(255);not null;default:'';index:idx_sku_key" comment:"sku_key"`
	IsConsumable            int     `json:"is_consumable" gorm:"type:tinyint;not null;default:0" comment:"是否是消耗品 0=否 1=是"`
	NoCodeNum               float64 `json:"no_code_num" gorm:"type:decimal(11,2);not null;default:0.00" comment:"无编码数量"`
	CodeNum                 float64 `json:"code_num" gorm:"type:decimal(11,2);not null;default:0.00" comment:"编码数量"`
	TotalNum                float64 `json:"total_num" gorm:"type:decimal(11,2);not null;default:0.00" comment:"总数量"`
	IsComponent             int     `json:"is_component" gorm:"type:tinyint(1);not null;default:0" comment:"是否是组成品，1 是 0 否"`

	// 关联关系 - 一对多
	StockProductDetails []StockProductDetail `json:"stock_product_details" gorm:"foreignkey:StockProductID"`
}

// ExistStockProductByID 根据ID检查产品是否存在
func ExistStockProductByID(id int) (bool, error) {
	var product StockProduct
	err := models.Db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if product.ID > 0 {
		return true, nil
	}

	return false, nil
}

// ExistStockProductByName 根据名称检查产品是否存在
func ExistStockProductByName(name string) (bool, error) {
	var product StockProduct
	err := models.Db.Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if product.ID > 0 {
		return true, nil
	}

	return false, nil
}

// GetStockProductTotal 获取产品总数
func GetStockProductTotal(maps interface{}) (int, error) {
	var count int
	if err := models.Db.Model(&StockProduct{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetStockProducts 获取产品列表（分页）
func GetStockProducts(pageNum int, pageSize int, maps interface{}) ([]*StockProduct, error) {
	var products []*StockProduct
	err := models.Db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return products, nil
}

// GetStockProducts1 获取产品列表（分页）- 版本2
func GetStockProducts1(pageNum int, pageSize int, maps interface{}) ([]*StockProduct, error) {
	var products []*StockProduct
	err := models.Db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return products, nil
}

// GetStockProduct 根据ID获取单个产品
func GetStockProduct(id int) (*StockProduct, error) {
	var product StockProduct
	err := models.Db.Where("id = ? AND deleted_on = ?", id, 0).First(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &product, nil
}

// GetStockProductWithDetails 根据ID获取产品及其明细
func GetStockProductWithDetails(id int) (*StockProduct, error) {
	var product StockProduct
	err := models.Db.Where("id = ? AND deleted_on = ?", id, 0).First(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 预加载关联的明细数据
	err = models.Db.Model(&product).Where("deleted_on = ?", 0).Related(&product.StockProductDetails, "StockProductID").Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &product, nil
}

// AddStockProduct 添加产品
func AddStockProduct(data map[string]interface{}) error {
	product := StockProduct{
		Unit:                    data["unit"].(string),
		Name:                    data["name"].(string),
		StockCustomizeProductID: data["stock_customize_product_id"].(int),
		SkuKey:                  data["sku_key"].(string),
		IsConsumable:            data["is_consumable"].(int),
		NoCodeNum:               data["no_code_num"].(float64),
		CodeNum:                 data["code_num"].(float64),
		TotalNum:                data["total_num"].(float64),
		IsComponent:             data["is_component"].(int),
	}

	if err := models.Db.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

// EditStockProduct 修改产品
func EditStockProduct(id int, data interface{}) error {
	if err := models.Db.Model(&StockProduct{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// DeleteStockProduct 删除产品（软删除）
func DeleteStockProduct(id int) error {
	if err := models.Db.Where("id = ?", id).Delete(StockProduct{}).Error; err != nil {
		return err
	}

	return nil
}

// UpdateStockProductNum 更新产品库存数量
func UpdateStockProductNum(id int, noCodeNum, codeNum float64) error {
	totalNum := noCodeNum + codeNum
	updates := map[string]interface{}{
		"no_code_num": noCodeNum,
		"code_num":    codeNum,
		"total_num":   totalNum,
	}

	if err := models.Db.Model(&StockProduct{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// GetStockProductsByCustomizeProductID 根据公司产品ID获取仓库产品列表
func GetStockProductsByCustomizeProductID(customizeProductID int) ([]*StockProduct, error) {
	var products []*StockProduct
	err := models.Db.Where("stock_customize_product_id = ? AND deleted_on = ?", customizeProductID, 0).Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return products, nil
}

// ===== GORM Hooks 示例 =====
// Hooks 是 GORM 提供的回调机制，在特定操作前后自动执行

// BeforeCreate 创建记录之前的回调
// 用途：数据验证、设置默认值、生成唯一标识等
func (p *StockProduct) BeforeCreate(scope *gorm.Scope) error {
	// 示例1：自动生成 SKU Key（如果未设置）
	if p.SkuKey == "" {
		p.SkuKey = generateSKU(p.Name)
	}

	// 示例2：数据验证
	if p.Name == "" {
		return gorm.ErrInvalidSQL
	}

	// 示例3：记录日志
	models.Db.Exec("INSERT INTO operation_logs (action, table_name, operation_time) VALUES (?, ?, NOW())",
		"CREATE", "stock_products")

	return nil
}

// AfterCreate 创建记录之后的回调
// 用途：触发其他业务逻辑、发送通知、记录审计日志等
func (p *StockProduct) AfterCreate(scope *gorm.Scope) error {
	// 示例1：记录创建日志
	models.Db.Exec("INSERT INTO audit_logs (entity_type, entity_id, action, user_id, created_at) VALUES (?, ?, ?, ?, NOW())",
		"stock_product", p.ID, "created", 1)

	// 示例2：发送通知（实际应该异步处理）
	// sendNotification("新产品创建", fmt.Sprintf("产品 %s 已创建", p.Name))

	return nil
}

// BeforeUpdate 更新记录之前的回调
// 用途：数据验证、记录变更历史、权限检查等
func (p *StockProduct) BeforeUpdate(scope *gorm.Scope) error {
	// 示例1：记录字段变更
	if scope.HasColumn("name") {
		oldName := ""
		scope.DB().Model(&StockProduct{}).Where("id = ?", p.ID).Select("name").Row().Scan(&oldName)
		if oldName != p.Name {
			// 记录名称变更历史
			models.Db.Exec("INSERT INTO change_history (entity_type, entity_id, field_name, old_value, new_value, changed_at) VALUES (?, ?, ?, ?, ?, NOW())",
				"stock_product", p.ID, "name", oldName, p.Name)
		}
	}

	// 示例2：更新修改时间
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}

// AfterUpdate 更新记录之后的回调
// 用途：同步缓存、触发其他系统更新等
func (p *StockProduct) AfterUpdate(scope *gorm.Scope) error {
	// 示例：清除相关缓存
	// cache.Delete(fmt.Sprintf("product:%d", p.ID))

	// 记录审计日志
	models.Db.Exec("INSERT INTO audit_logs (entity_type, entity_id, action, user_id, created_at) VALUES (?, ?, ?, ?, NOW())",
		"stock_product", p.ID, "updated", 1)

	return nil
}

// BeforeDelete 删除记录之前的回调
// 用途：权限检查、级联删除检查、备份数据等
func (p *StockProduct) BeforeDelete(scope *gorm.Scope) error {
	// 示例1：检查是否有关联的明细记录
	var count int
	models.Db.Model(&StockProductDetail{}).Where("stock_product_id = ?", p.ID).Count(&count)
	if count > 0 {
		// 如果有明细记录，可以选择：
		// 1. 阻止删除
		// return errors.New("产品下还有明细记录，不能删除")
		// 2. 或者级联软删除明细记录
		models.Db.Model(&StockProductDetail{}).Where("stock_product_id = ?", p.ID).Update("deleted_on", time.Now().Unix())
	}

	// 示例2：备份要删除的数据
	models.Db.Exec("INSERT INTO deleted_products_backup SELECT * FROM stock_products WHERE id = ?", p.ID)

	return nil
}

// AfterDelete 删除记录之后的回调
// 用途：清理关联数据、发送通知等
func (p *StockProduct) AfterDelete(scope *gorm.Scope) error {
	// 记录删除日志
	models.Db.Exec("INSERT INTO audit_logs (entity_type, entity_id, action, user_id, created_at) VALUES (?, ?, ?, ?, NOW())",
		"stock_product", p.ID, "deleted", 1)

	return nil
}

// BeforeSave 保存（创建或更新）之前的回调
// 用途：通用的数据处理，不区分创建还是更新
func (p *StockProduct) BeforeSave(scope *gorm.Scope) error {
	// 示例：自动计算总数量
	p.TotalNum = p.NoCodeNum + p.CodeNum

	return nil
}

// AfterFind 查询记录之后的回调
// 用途：数据解密、格式化、加载额外信息等
func (p *StockProduct) AfterFind(scope *gorm.Scope) error {
	// 示例：加载额外的统计信息
	// 注意：这个操作可能会导致 N+1 查询问题，实际应用中要谨慎使用
	// 或者考虑使用 Preload

	return nil
}

// generateSKU 生成 SKU Key 的辅助函数
func generateSKU(name string) string {
	// 简单示例：使用时间戳和名称生成
	return fmt.Sprintf("SKU-%s-%d", name, time.Now().Unix())
}
