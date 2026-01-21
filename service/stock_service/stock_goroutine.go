package stock_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/models/stock"
)

// ===== 协程相关类型定义 =====

type StockUpdate struct {
	ProductID int     `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

type UpdateResult struct {
	ProductID int    `json:"product_id"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

type ProductResult struct {
	ProductID int                    `json:"product_id"`
	Data      map[string]interface{} `json:"data"`
	Error     string                 `json:"error,omitempty"`
}

// ===== 协程并发处理示例 =====

// BatchQueryProductsWithGoroutine 并发查询多个产品（演示基本协程使用）
func BatchQueryProductsWithGoroutine(productIDs []int) ([]ProductResult, error) {
	results := make([]ProductResult, len(productIDs))
	var wg sync.WaitGroup
	var mu sync.Mutex // 保护 results 切片

	// 为每个产品ID启动一个协程
	for i, productID := range productIDs {
		wg.Add(1)
		go func(index, id int) {
			defer wg.Done()

			// 查询产品信息
			var product stock.StockProduct
			if err := models.Db.Where("id = ? AND deleted_on = ?", id, 0).First(&product).Error; err != nil {
				mu.Lock()
				results[index] = ProductResult{
					ProductID: id,
					Error:     fmt.Sprintf("查询失败: %v", err),
				}
				mu.Unlock()
				return
			}

			// 查询产品明细
			var details []stock.StockProductDetail
			models.Db.Where("stock_product_id = ? AND deleted_on = ?", id, 0).Find(&details)

			// 构建结果
			mu.Lock()
			results[index] = ProductResult{
				ProductID: id,
				Data: map[string]interface{}{
					"id":           product.ID,
					"name":         product.Name,
					"total_num":    product.TotalNum,
					"unit":         product.Unit,
					"detail_count": len(details),
				},
			}
			mu.Unlock()
		}(i, productID)
	}

	// 等待所有协程完成
	wg.Wait()

	return results, nil
}

// BatchUpdateStockWithGoroutine 并发更新多个产品库存（演示带超时的协程）
func BatchUpdateStockWithGoroutine(updates []StockUpdate) ([]UpdateResult, error) {
	results := make([]UpdateResult, len(updates))
	var wg sync.WaitGroup

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用 channel 收集结果
	resultChan := make(chan UpdateResult, len(updates))
	errorChan := make(chan error, len(updates))

	for i, update := range updates {
		wg.Add(1)
		go func(index int, u StockUpdate) {
			defer wg.Done()

			// 检查 context 是否已取消
			select {
			case <-ctx.Done():
				resultChan <- UpdateResult{
					ProductID: u.ProductID,
					Success:   false,
					Error:     "超时",
				}
				return
			default:
			}

			// 模拟耗时操作
			time.Sleep(100 * time.Millisecond)

			// 更新库存
			err := models.Db.Model(&stock.StockProduct{}).
				Where("id = ? AND deleted_on = ?", u.ProductID, 0).
				Update("total_num", u.Quantity).Error

			if err != nil {
				resultChan <- UpdateResult{
					ProductID: u.ProductID,
					Success:   false,
					Error:     err.Error(),
				}
			} else {
				resultChan <- UpdateResult{
					ProductID: u.ProductID,
					Success:   true,
				}
			}
		}(i, update)
	}

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// 收集结果
	i := 0
	for result := range resultChan {
		results[i] = result
		i++
	}

	return results, nil
}

// ProcessWithWorkerPool Worker Pool 模式（演示固定数量的工作协程）
func ProcessWithWorkerPool(productIDs []int, workerCount int) ([]ProductResult, error) {
	// 创建任务 channel 和结果 channel
	jobChan := make(chan int, len(productIDs))
	resultChan := make(chan ProductResult, len(productIDs))

	// 创建 worker pool
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i+1, jobChan, resultChan, &wg)
	}

	// 发送任务到 jobChan
	go func() {
		for _, id := range productIDs {
			jobChan <- id
		}
		close(jobChan) // 关闭任务通道，告诉 worker 没有更多任务了
	}()

	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(resultChan) // 所有 worker 完成后关闭结果通道
	}()

	// 收集结果
	var results []ProductResult
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}

// worker Worker Pool 中的工作协程
func worker(id int, jobs <-chan int, results chan<- ProductResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for productID := range jobs {
		// 模拟处理任务
		time.Sleep(50 * time.Millisecond)

		var product stock.StockProduct
		err := models.Db.Where("id = ? AND deleted_on = ?", productID, 0).First(&product).Error

		if err != nil {
			results <- ProductResult{
				ProductID: productID,
				Error:     fmt.Sprintf("Worker %d 处理失败: %v", id, err),
			}
			continue
		}

		results <- ProductResult{
			ProductID: productID,
			Data: map[string]interface{}{
				"id":        product.ID,
				"name":      product.Name,
				"total_num": product.TotalNum,
				"worker_id": id,
			},
		}
	}
}

// ProcessWithPipeline Pipeline 模式（演示数据流处理）
func ProcessWithPipeline(limit int) ([]map[string]interface{}, error) {
	// Stage 1: 生成产品 ID
	idChan := generateProductIDs(limit)

	// Stage 2: 获取产品信息
	productChan := fetchProducts(idChan)

	// Stage 3: 计算统计信息
	resultChan := calculateStats(productChan)

	// 收集最终结果
	var results []map[string]interface{}
	for result := range resultChan {
		results = append(results, result)
	}

	return results, nil
}

// generateProductIDs Pipeline 的第一阶段：生成产品ID
func generateProductIDs(limit int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)

		var products []stock.StockProduct
		models.Db.Where("deleted_on = ?", 0).Limit(limit).Find(&products)

		for _, product := range products {
			out <- product.ID
		}
	}()
	return out
}

// fetchProducts Pipeline 的第二阶段：获取产品信息
func fetchProducts(ids <-chan int) <-chan stock.StockProduct {
	out := make(chan stock.StockProduct)
	go func() {
		defer close(out)

		for id := range ids {
			var product stock.StockProduct
			if err := models.Db.Where("id = ?", id).First(&product).Error; err == nil {
				out <- product
			}
		}
	}()
	return out
}

// calculateStats Pipeline 的第三阶段：计算统计信息
func calculateStats(products <-chan stock.StockProduct) <-chan map[string]interface{} {
	out := make(chan map[string]interface{})
	go func() {
		defer close(out)

		for product := range products {
			// 查询该产品的明细数量
			var count int
			models.Db.Model(&stock.StockProductDetail{}).
				Where("stock_product_id = ? AND deleted_on = ?", product.ID, 0).
				Count(&count)

			out <- map[string]interface{}{
				"product_id":   product.ID,
				"product_name": product.Name,
				"total_num":    product.TotalNum,
				"detail_count": count,
				"avg_per_detail": func() float64 {
					if count > 0 {
						return product.TotalNum / float64(count)
					}
					return 0
				}(),
			}
		}
	}()
	return out
}

// ProcessWithFanOutFanIn Fan-out/Fan-in 模式（演示多个协程处理，单个协程汇总）
func ProcessWithFanOutFanIn(productIDs []int) (map[string]interface{}, error) {
	// Fan-out: 将数据分发到多个协程处理
	numWorkers := 3
	inputChan := make(chan int, len(productIDs))

	// 发送所有产品ID到输入通道
	go func() {
		for _, id := range productIDs {
			inputChan <- id
		}
		close(inputChan)
	}()

	// 创建多个输出通道（每个 worker 一个）
	channels := make([]<-chan ProductResult, numWorkers)
	for i := 0; i < numWorkers; i++ {
		channels[i] = processProducts(inputChan)
	}

	// Fan-in: 合并多个通道的结果
	results := fanIn(channels...)

	// 收集并统计结果
	var allResults []ProductResult
	totalNum := 0.0
	successCount := 0

	for result := range results {
		allResults = append(allResults, result)
		if result.Error == "" {
			successCount++
			if num, ok := result.Data["total_num"].(float64); ok {
				totalNum += num
			}
		}
	}

	return map[string]interface{}{
		"total_products":  len(allResults),
		"success_count":   successCount,
		"total_stock_num": totalNum,
		"average_per_item": func() float64 {
			if successCount > 0 {
				return totalNum / float64(successCount)
			}
			return 0
		}(),
		"details": allResults,
	}, nil
}

// processProducts Fan-out 阶段：处理产品
func processProducts(input <-chan int) <-chan ProductResult {
	out := make(chan ProductResult)
	go func() {
		defer close(out)

		for id := range input {
			var product stock.StockProduct
			err := models.Db.Where("id = ? AND deleted_on = ?", id, 0).First(&product).Error

			if err != nil {
				out <- ProductResult{
					ProductID: id,
					Error:     err.Error(),
				}
				continue
			}

			// 模拟耗时处理
			time.Sleep(10 * time.Millisecond)

			out <- ProductResult{
				ProductID: id,
				Data: map[string]interface{}{
					"id":        product.ID,
					"name":      product.Name,
					"total_num": product.TotalNum,
				},
			}
		}
	}()
	return out
}

// fanIn Fan-in 阶段：合并多个通道
func fanIn(channels ...<-chan ProductResult) <-chan ProductResult {
	out := make(chan ProductResult)
	var wg sync.WaitGroup

	// 为每个输入通道启动一个协程
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan ProductResult) {
			defer wg.Done()
			for result := range c {
				out <- result
			}
		}(ch)
	}

	// 等待所有协程完成后关闭输出通道
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
