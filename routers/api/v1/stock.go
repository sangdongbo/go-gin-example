package v1

import (
	"net/http"
	"strconv"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/stock_service"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// ===== 仓库产品相关接口 =====

// @Summary 获取仓库产品列表
// @Tags 仓库产品管理
// @Produce json
// @Param name query string false "产品名称"
// @Param stock_product_id query int false "公司产品ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/products [get]
func GetStockProducts(c *gin.Context) {
	appG := app.Gin{C: c}

	customizeProductID := com.StrTo(c.Query("stock_product_id")).MustInt()

	query := stock_service.StockProductQuery{
		Name:                    c.Query("name"),
		StockCustomizeProductID: customizeProductID,
		PageNum:                 util.GetPage(c),
		PageSize:                setting.AppSetting.PageSize,
	}

	products, err := query.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	count, err := query.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": products,
		"total": count,
	})
}

func GetStockProducts1(c *gin.Context) {
	appG := app.Gin{C: c}

	customizeProductID := com.StrTo(c.Query("stock_product_id")).MustInt()

	query := stock_service.StockProductQuery{
		Name:                    c.Query("name"),
		StockCustomizeProductID: customizeProductID,
		PageNum:                 util.GetPage(c),
		PageSize:                setting.AppSetting.PageSize,
	}

	products, err := query.GetAll1()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	count, err := query.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": products,
		"total": count,
	})
}

// @Summary 获取单个仓库产品
// @Tags 仓库产品管理
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product/{id} [get]
func GetStockProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	id, _ := strconv.Atoi(c.Param("id"))

	exists, err := stock_service.ExistStockProductByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusNotFound, e.ERROR_NOT_EXIST, nil)
		return
	}

	product, err := stock_service.GetStockProductByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}

// @Summary 获取仓库产品及明细
// @Tags 仓库产品管理
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product/{id}/details [get]
func GetStockProductWithDetails(c *gin.Context) {
	appG := app.Gin{C: c}

	id, _ := strconv.Atoi(c.Param("id"))

	product, err := stock_service.GetStockProductWithDetails(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}

// @Summary 创建仓库产品
// @Tags 仓库产品管理
// @Accept json
// @Produce json
// @Param product body stock_service.AddStockProductForm true "产品信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product [post]
func AddStockProduct(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form stock_service.AddStockProductForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 检查产品名称是否已存在
	exists, err := stock_service.ExistStockProductByName(form.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST, nil)
		return
	}

	product := stock_service.ConvertAddFormToStockProduct(form)
	if err := product.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 更新仓库产品
// @Tags 仓库产品管理
// @Accept json
// @Produce json
// @Param product body stock_service.EditStockProductForm true "产品信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product [put]
func UpdateStockProduct(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form stock_service.EditStockProductForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 验证产品是否存在
	exists, err := stock_service.ExistStockProductByID(form.ID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusNotFound, e.ERROR_NOT_EXIST, nil)
		return
	}

	product := stock_service.ConvertEditFormToStockProduct(form)
	if err := product.Edit(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除仓库产品
// @Tags 仓库产品管理
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product/{id} [delete]
func DeleteStockProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	id, _ := strconv.Atoi(c.Param("id"))

	deleteService := stock_service.StockProductDelete{ID: id}
	err := deleteService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除成功")
}

// ===== 产品明细相关接口 =====

// @Summary 获取产品明细列表
// @Tags 产品明细管理
// @Produce json
// @Param stock_product_id query int false "产品ID"
// @Param order_id query int false "订单ID"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product-details [get]
func GetStockProductDetails(c *gin.Context) {
	appG := app.Gin{C: c}

	productID := com.StrTo(c.Query("stock_product_id")).MustInt()
	orderID := com.StrTo(c.Query("order_id")).MustInt()

	status := -1
	if c.Query("status") != "" {
		status = com.StrTo(c.Query("status")).MustInt()
	}

	query := stock_service.StockProductDetailQuery{
		StockProductID: productID,
		OrderID:        orderID,
		Status:         status,
		PageNum:        util.GetPage(c),
		PageSize:       setting.AppSetting.PageSize,
	}

	details, err := query.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	count, err := query.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": details,
		"total": count,
	})
}

// @Summary 获取单个产品明细
// @Tags 产品明细管理
// @Produce json
// @Param id path int true "明细ID"
// @Success 200 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product-detail/{id} [get]
func GetStockProductDetail(c *gin.Context) {
	appG := app.Gin{C: c}

	id, _ := strconv.Atoi(c.Param("id"))

	exists, err := stock_service.ExistStockProductDetailByID(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusNotFound, e.ERROR_NOT_EXIST, nil)
		return
	}

	detail, err := stock_service.GetStockProductDetailWithProduct(id)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, detail)
}

// @Summary 创建产品明细
// @Tags 产品明细管理
// @Accept json
// @Produce json
// @Param detail body stock_service.AddStockProductDetailForm true "明细信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product-detail [post]
func AddStockProductDetail(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form stock_service.AddStockProductDetailForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 检查唯一编码是否已存在
	exists, err := stock_service.ExistStockProductDetailByHiddenCode(form.HiddenCode)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST, nil)
		return
	}

	detail := stock_service.ConvertAddFormToStockProductDetail(form)
	if err := detail.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 更新产品明细
// @Tags 产品明细管理
// @Accept json
// @Produce json
// @Param detail body stock_service.EditStockProductDetailForm true "明细信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product-detail [put]
func UpdateStockProductDetail(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form stock_service.EditStockProductDetailForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	// 验证明细是否存在
	exists, err := stock_service.ExistStockProductDetailByID(form.ID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusNotFound, e.ERROR_NOT_EXIST, nil)
		return
	}

	detail := stock_service.ConvertEditFormToStockProductDetail(form)
	if err := detail.Edit(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除产品明细
// @Tags 产品明细管理
// @Produce json
// @Param id path int true "明细ID"
// @Success 200 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product-detail/{id} [delete]
func DeleteStockProductDetail(c *gin.Context) {
	appG := app.Gin{C: c}

	id, _ := strconv.Atoi(c.Param("id"))

	deleteService := stock_service.StockProductDetailDelete{ID: id}
	err := deleteService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除成功")
}

// @Summary 获取产品明细汇总
// @Tags 产品明细管理
// @Produce json
// @Param product_id query int true "产品ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/product-detail/summary [get]
func GetStockProductDetailSummary(c *gin.Context) {
	appG := app.Gin{C: c}

	productID, _ := strconv.Atoi(c.Query("product_id"))

	summary, err := stock_service.GetStockProductDetailSummary(productID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, summary)
}

// ===== 事务相关接口 =====

// @Summary 批量创建产品及明细（事务）
// @Tags 仓库产品管理-事务
// @Produce json
// @Param body body object true "批量创建数据"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/transaction/batch-create [post]
func BatchCreateProductWithDetails(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		Products []struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			Details     []struct {
				Quantity    int     `json:"quantity" binding:"required"`
				Price       float64 `json:"price" binding:"required"`
				Location    string  `json:"location"`
				Supplier    string  `json:"supplier"`
				BatchNumber string  `json:"batch_number"`
			} `json:"details" binding:"required,min=1"`
		} `json:"products" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 转换为 service 层的结构体
	var products []stock_service.BatchProduct
	for _, p := range req.Products {
		var details []stock_service.BatchProductDetail
		for _, d := range p.Details {
			details = append(details, stock_service.BatchProductDetail{
				Quantity:    d.Quantity,
				Price:       d.Price,
				Location:    d.Location,
				Supplier:    d.Supplier,
				BatchNumber: d.BatchNumber,
			})
		}
		products = append(products, stock_service.BatchProduct{
			Name:        p.Name,
			Description: p.Description,
			Details:     details,
		})
	}

	if err := stock_service.BatchCreateProductWithDetails(products); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "批量创建成功")
}

// @Summary 库存转移（事务）
// @Tags 仓库产品管理-事务
// @Produce json
// @Param body body object true "转移数据"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/transaction/transfer [post]
func TransferStock(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		FromDetailID int `json:"from_detail_id" binding:"required"`
		ToDetailID   int `json:"to_detail_id" binding:"required"`
		Quantity     int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := stock_service.TransferStock(req.FromDetailID, req.ToDetailID, req.Quantity); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "库存转移成功")
}

// ===== Hooks 相关接口 =====

// @Summary 创建产品（演示 Hooks）
// @Tags 仓库产品管理-Hooks
// @Produce json
// @Param body body object true "产品数据"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/hooks/product [post]
func CreateProductWithHooks(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	product, err := stock_service.CreateProductWithHooks(req.Name, req.Description)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}

// @Summary 更新产品（演示 Hooks）
// @Tags 仓库产品管理-Hooks
// @Produce json
// @Param id path int true "产品ID"
// @Param body body object true "产品数据"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/hooks/product/{id} [put]
func UpdateProductWithHooks(c *gin.Context) {
	appG := app.Gin{C: c}

	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := stock_service.UpdateProductWithHooks(id, req.Name, req.Description); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "更新成功")
}

// ===== 关联查询优化 =====

// @Summary 预加载关联数据（Preload 优化）
// @Tags 仓库产品管理-关联查询优化
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/optimize/products-with-details [get]
func GetProductsWithDetailsOptimized(c *gin.Context) {
	appG := app.Gin{C: c}

	pageNum := util.GetPage(c)
	pageSize := setting.AppSetting.PageSize

	products, total, err := stock_service.GetProductsWithDetailsOptimized(pageNum, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": products,
		"total": total,
	})
}

// @Summary 使用 Join 查询（Join 优化）
// @Tags 仓库产品管理-关联查询优化
// @Produce json
// @Param min_quantity query int false "最小库存数量"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/optimize/products-join [get]
func GetProductsWithJoin(c *gin.Context) {
	appG := app.Gin{C: c}

	minQuantity := com.StrTo(c.Query("min_quantity")).MustInt()

	products, err := stock_service.GetProductsWithJoin(minQuantity)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, products)
}

// ===== 事务嵌套 =====

// @Summary 创建订单（嵌套事务）
// @Tags 仓库产品管理-事务嵌套
// @Produce json
// @Param body body object true "订单数据"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/nested-transaction/order [post]
func CreateOrderWithNestedTransaction(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		OrderSN string `json:"order_sn" binding:"required"`
		Items   []struct {
			ProductID int `json:"product_id" binding:"required"`
			Quantity  int `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 转换为 service 层的结构体
	var items []stock_service.OrderItem
	for _, item := range req.Items {
		items = append(items, stock_service.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	if err := stock_service.CreateOrderWithNestedTransaction(req.OrderSN, items); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "订单创建成功")
}

// ===== 协程并发处理 =====

// @Summary 并发查询多个产品
// @Tags 仓库产品管理-协程
// @Produce json
// @Param body body object true "产品ID列表"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/goroutine/batch-query [post]
func BatchQueryProductsWithGoroutine(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		ProductIDs []int `json:"product_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	products, err := stock_service.BatchQueryProductsWithGoroutine(req.ProductIDs)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, products)
}

// @Summary 并发更新库存
// @Tags 仓库产品管理-协程
// @Produce json
// @Param body body object true "更新数据"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/goroutine/batch-update [post]
func BatchUpdateStockWithGoroutine(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		Updates []struct {
			ProductID int     `json:"product_id" binding:"required"`
			Quantity  float64 `json:"quantity" binding:"required"`
		} `json:"updates" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	var updates []stock_service.StockUpdate
	for _, u := range req.Updates {
		updates = append(updates, stock_service.StockUpdate{
			ProductID: u.ProductID,
			Quantity:  u.Quantity,
		})
	}

	results, err := stock_service.BatchUpdateStockWithGoroutine(updates)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, results)
}

// @Summary Worker Pool 模式处理任务
// @Tags 仓库产品管理-协程
// @Produce json
// @Param body body object true "任务列表"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/goroutine/worker-pool [post]
func ProcessWithWorkerPool(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		ProductIDs  []int `json:"product_ids" binding:"required,min=1"`
		WorkerCount int   `json:"worker_count" binding:"min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 默认 worker 数量
	if req.WorkerCount == 0 {
		req.WorkerCount = 5
	}

	results, err := stock_service.ProcessWithWorkerPool(req.ProductIDs, req.WorkerCount)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, results)
}

// @Summary Pipeline 模式处理数据
// @Tags 仓库产品管理-协程
// @Produce json
// @Param limit query int false "数据限制"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/goroutine/pipeline [get]
func ProcessWithPipeline(c *gin.Context) {
	appG := app.Gin{C: c}

	limit := com.StrTo(c.Query("limit")).MustInt()
	if limit == 0 {
		limit = 10
	}

	results, err := stock_service.ProcessWithPipeline(limit)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, results)
}

// @Summary Fan-out/Fan-in 模式
// @Tags 仓库产品管理-协程
// @Produce json
// @Param body body object true "产品ID列表"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/stock/goroutine/fan-out-in [post]
func ProcessWithFanOutFanIn(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		ProductIDs []int `json:"product_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	results, err := stock_service.ProcessWithFanOutFanIn(req.ProductIDs)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, results)
}
