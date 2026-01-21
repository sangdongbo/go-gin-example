package stock_service

import (
	"fmt"
	"time"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/models/stock"
	"github.com/jinzhu/gorm"
)

// ===== 事务相关类型定义 =====

type BatchProductDetail struct {
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Location    string  `json:"location"`
	Supplier    string  `json:"supplier"`
	BatchNumber string  `json:"batch_number"`
}

type BatchProduct struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Details     []BatchProductDetail `json:"details"`
}

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// ===== 事务相关 =====

// BatchCreateProductWithDetails 批量创建产品及明细（演示事务使用）
func BatchCreateProductWithDetails(productList []BatchProduct) error {
	// 开始事务
	tx := models.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 批量创建产品和明细
	for _, p := range productList {
		// 创建产品
		product := stock.StockProduct{
			Name:                    p.Name,
			Unit:                    "个",
			StockCustomizeProductID: 1,
			IsConsumable:            0,
			IsComponent:             0,
		}

		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建产品失败: %v", err)
		}

		// 创建产品明细
		for _, d := range p.Details {
			detail := stock.StockProductDetail{
				StockProductID: product.ID,
				Num:            float64(d.Quantity),
				CostPrice:      d.Price,
				NeedReturn:     0,
				Status:         0,
				Note:           fmt.Sprintf("Location: %s, Supplier: %s, Batch: %s", d.Location, d.Supplier, d.BatchNumber),
			}

			if err := tx.Create(&detail).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("创建产品明细失败: %v", err)
			}

			// 更新产品总数量
			product.TotalNum += float64(d.Quantity)
		}

		// 更新产品的总数量
		if err := tx.Model(&product).Update("total_num", product.TotalNum).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新产品总数量失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// TransferStock 库存转移（演示事务的 ACID 特性）
func TransferStock(fromDetailID, toDetailID, quantity int) error {
	// 开始事务
	tx := models.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 查询源库存明细
	var fromDetail stock.StockProductDetail
	if err := tx.Where("id = ?", fromDetailID).First(&fromDetail).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("源库存不存在: %v", err)
	}

	// 查询目标库存明细
	var toDetail stock.StockProductDetail
	if err := tx.Where("id = ?", toDetailID).First(&toDetail).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("目标库存不存在: %v", err)
	}

	// 验证库存是否充足
	if int(fromDetail.Num) < quantity {
		tx.Rollback()
		return fmt.Errorf("源库存不足，当前库存: %.2f, 需要转移: %d", fromDetail.Num, quantity)
	}

	// 减少源库存
	if err := tx.Model(&fromDetail).Update("num", fromDetail.Num-float64(quantity)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("减少源库存失败: %v", err)
	}

	// 增加目标库存
	if err := tx.Model(&toDetail).Update("num", toDetail.Num+float64(quantity)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("增加目标库存失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// ===== Hooks 相关 =====

// CreateProductWithHooks 创建产品（演示 GORM Hooks）
// 注意：Hooks 需要在 Model 中定义，这里只是调用创建
func CreateProductWithHooks(name, description string) (*stock.StockProduct, error) {
	product := stock.StockProduct{
		Name:                    name,
		Unit:                    "个",
		StockCustomizeProductID: 1,
		IsConsumable:            0,
		IsComponent:             0,
	}

	if err := models.Db.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("创建产品失败: %v", err)
	}

	return &product, nil
}

// UpdateProductWithHooks 更新产品（演示 GORM Hooks）
func UpdateProductWithHooks(id int, name, description string) error {
	updates := map[string]interface{}{}

	if name != "" {
		updates["name"] = name
	}

	if err := models.Db.Model(&stock.StockProduct{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新产品失败: %v", err)
	}

	return nil
}

// ===== 关联查询优化 =====

// GetProductsWithDetailsOptimized 使用 Preload 预加载关联数据（避免 N+1 查询问题）
func GetProductsWithDetailsOptimized(pageNum, pageSize int) ([]stock.StockProduct, int, error) {
	var products []stock.StockProduct
	var total int

	// 先统计总数
	if err := models.Db.Model(&stock.StockProduct{}).Where("deleted_on = ?", 0).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 使用 Preload 预加载关联数据，一次性加载产品和明细
	// 避免了 N+1 查询问题：不使用 Preload 时，查询 N 个产品需要 1 + N 次查询
	// 使用 Preload 后，只需要 2 次查询（1次查产品，1次查所有明细）
	offset := (pageNum - 1) * pageSize
	if err := models.Db.Preload("StockProductDetails").
		Where("deleted_on = ?", 0).
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductsWithJoin 使用 Joins 优化查询（适合需要关联条件过滤的场景）
func GetProductsWithJoin(minQuantity int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 使用 Joins 进行关联查询，只返回库存数量大于指定值的产品
	// Joins 会生成 INNER JOIN SQL，适合需要根据关联表条件筛选的场景
	rows, err := models.Db.Table("stock_products").
		Select("stock_products.id, stock_products.name, stock_products.total_num, "+
			"stock_product_details.quantity, stock_product_details.location, stock_product_details.supplier").
		Joins("INNER JOIN stock_product_details ON stock_products.id = stock_product_details.stock_product_id").
		Where("stock_products.deleted_on = ? AND stock_product_details.quantity >= ?", 0, minQuantity).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var totalNum float64
		var quantity int
		var location, supplier string

		if err := rows.Scan(&id, &name, &totalNum, &quantity, &location, &supplier); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":        id,
			"name":      name,
			"total_num": totalNum,
			"detail": map[string]interface{}{
				"quantity": quantity,
				"location": location,
				"supplier": supplier,
			},
		})
	}

	return results, nil
}

// ===== 事务嵌套 =====

// CreateOrderWithNestedTransaction 创建订单（演示嵌套事务）
func CreateOrderWithNestedTransaction(orderSN string, orderItems []OrderItem) error {
	// 外层事务：处理订单创建
	tx := models.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// 创建订单（这里简化处理，实际应该有专门的订单表）
	orderData := map[string]interface{}{
		"order_sn":   orderSN,
		"status":     1,
		"created_at": time.Now(),
	}

	// 注意：这里演示嵌套事务的概念
	// 在同一个事务中处理多个相关操作
	for _, item := range orderItems {
		// 内层操作1：检查产品是否存在
		var product stock.StockProduct
		if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("产品不存在: %v", err)
		}

		// 内层操作2：检查库存是否充足（从明细表中查询）
		var totalQuantity float64
		if err := tx.Model(&stock.StockProductDetail{}).
			Where("stock_product_id = ?", item.ProductID).
			Select("COALESCE(SUM(num), 0)").
			Row().Scan(&totalQuantity); err != nil {
			tx.Rollback()
			return fmt.Errorf("查询库存失败: %v", err)
		}

		if int(totalQuantity) < item.Quantity {
			tx.Rollback()
			return fmt.Errorf("产品 %s 库存不足，需要: %d, 可用: %.2f",
				product.Name, item.Quantity, totalQuantity)
		}

		// 内层操作3：扣减库存（使用 SavePoint 实现真正的嵌套事务）
		// 注意：MySQL 的嵌套事务实际上是通过 SAVEPOINT 实现的
		tx.Exec("SAVEPOINT sp_reduce_stock")

		if err := reduceStockInTransaction(tx, item.ProductID, item.Quantity); err != nil {
			// 如果扣减失败，回滚到 SAVEPOINT
			tx.Exec("ROLLBACK TO SAVEPOINT sp_reduce_stock")
			tx.Rollback()
			return fmt.Errorf("扣减库存失败: %v", err)
		}

		// 释放 SAVEPOINT
		tx.Exec("RELEASE SAVEPOINT sp_reduce_stock")
	}

	// 记录订单信息到日志或其他表
	fmt.Printf("订单创建成功: %+v\n", orderData)

	// 提交外层事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// reduceStockInTransaction 在事务中扣减库存（内层事务函数）
func reduceStockInTransaction(tx *gorm.DB, productID, quantity int) error {
	// 查询可用库存明细（按创建时间排序，先进先出）
	var details []stock.StockProductDetail
	if err := tx.Where("stock_product_id = ? AND num > 0", productID).
		Order("created_on ASC").
		Find(&details).Error; err != nil {
		return err
	}

	remainingQuantity := quantity

	for _, detail := range details {
		if remainingQuantity <= 0 {
			break
		}

		// 计算本次扣减数量
		deductQuantity := remainingQuantity
		if int(detail.Num) < remainingQuantity {
			deductQuantity = int(detail.Num)
		}

		// 更新库存明细
		newNum := detail.Num - float64(deductQuantity)
		if err := tx.Model(&detail).Update("num", newNum).Error; err != nil {
			return err
		}

		remainingQuantity -= deductQuantity
	}

	if remainingQuantity > 0 {
		return fmt.Errorf("库存不足，还需要 %d 个", remainingQuantity)
	}

	return nil
}
