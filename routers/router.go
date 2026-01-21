package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/EDDYCJY/go-gin-example/docs"
	"github.com/EDDYCJY/go-gin-example/middleware/jwt"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/EDDYCJY/go-gin-example/pkg/export"
	"github.com/EDDYCJY/go-gin-example/pkg/qrcode"
	"github.com/EDDYCJY/go-gin-example/pkg/upload"
	"github.com/EDDYCJY/go-gin-example/routers/api"
	v1 "github.com/EDDYCJY/go-gin-example/routers/api/v1"
	v1Test "github.com/EDDYCJY/go-gin-example/routers/api/v1/test"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	// 协程测试（不需要登录）
	coroutine := r.Group("/api/v1/coroutine")
	{
		coroutine.POST("test", v1.TestCoroutine) // 测试协程处理
	}

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// Casbin 权限管理接口
		casbin := apiv1.Group("/casbin")
		{
			// 角色管理
			casbin.POST("/create-role", v1.CreateRole)           // 创建新角色
			casbin.DELETE("/delete-role-by-name", v1.DeleteRole) // 删除角色

			casbin.GET("/roles", v1.GetAllRoles)                      // 获取所有角色列表
			casbin.GET("/role-permissions", v1.GetPermissionsForRole) // 获取角色的所有权限
			casbin.GET("/role-users", v1.GetUsersForRole)             // 获取拥有某个角色的所有用户

			// 用户角色管理
			casbin.POST("/add-role", v1.AddRoleForUser)         // 为用户添加角色
			casbin.DELETE("/delete-role", v1.DeleteRoleForUser) // 删除用户角色
			casbin.GET("/user-roles", v1.GetRolesForUser)       // 获取用户的所有角色

			// 权限策略管理
			casbin.POST("/add-policy", v1.AddPolicy)            // 添加权限策略
			casbin.DELETE("/delete-policy", v1.DeletePolicy)    // 删除权限策略
			casbin.GET("/check-permission", v1.CheckPermission) // 检查权限
		}

		// 添加用户
		amqp := apiv1.Group("/amqp")
		{
			amqp.POST("addMqUser", v1.AddMqUsers) //添加MQ用户
			amqp.GET("consume", v1.ConsumeMessage)
			amqp.GET("consumeAck", v1.ConsumeAckMessage)
			amqp.GET("sendDeadlineMessage", v1.SendDeadlineMessage)
			amqp.GET("consumeDeadlineMessageOne", v1.ConsumeDeadlineMessageOne)
			amqp.GET("consumeDeadlineMessageTwo", v1.ConsumeDeadlineMessageTwo)

			amqp.POST("addEmail", v1.AddEmails)
			amqp.PUT("updateEmail", v1.UpdateEmail)
			amqp.GET("getEmails", v1.GetEmails)
		}

		es := apiv1.Group("/es")
		{
			es.POST("addEsData", v1.AddEsData) //添加ES 测试数据
			es.GET("getEsData", v1.GetEsData)

			es.POST("addEsJsonData", v1.AddEsJsonData)
			es.GET("searchProductByKeyword", v1.SearchProductByKeyword)
			es.GET("searchByBrandOrigin", v1.SearchByBrandOrigin)
			es.GET("filterProductByPrice", v1.FilterProductByPrice)
		}

		gmp := apiv1.Group("/gmp")
		{
			gmp.GET("testOneContext", v1.TestOneContext)
			gmp.GET("limitRate", v1.LimitRate)
			gmp.GET("limitRate1", v1.LimitRate1)
			gmp.GET("aggregate", v1.Aggregate)
			gmp.GET("fastest", v1.Fastest)
			gmp.GET("testContextTimeout", v1.TestContextTimeout)
			gmp.GET("aggregate2", v1.Aggregate2)
			gmp.GET("doTask1", v1.DoTask1)
		}

		test := apiv1.Group("/test")
		{
			test.GET("testMap", v1Test.TestMap)
			test.GET("testArray", v1Test.TestArray)
			test.GET("testSlice", v1Test.TestSlice)
			test.GET("goroutine", v1Test.TestGoroutine)
			test.GET("testStruct", v1Test.TestStruct)
		}

		// tag 分组
		orders := apiv1.Group("/orders")
		{
			orders.GET("", v1.GetOrders)            // 获取订单列表
			orders.POST("addOrder", v1.AddOrder)    //添加订单
			orders.PUT("editOrder", v1.UpdateOrder) //修改订单
			orders.DELETE("deleteOrder/:order_sn", v1.DeleteOrder)
		}

		// 仓库产品管理
		stock := apiv1.Group("/stock")
		{
			// 产品相关
			stock.GET("/products", v1.GetStockProducts) // 获取产品列表
			stock.GET("/products1", v1.GetStockProducts1)
			stock.GET("/product/:id", v1.GetStockProduct)                    // 获取单个产品
			stock.GET("/product/:id/details", v1.GetStockProductWithDetails) // 获取产品及明细
			stock.POST("/product", v1.AddStockProduct)                       // 创建产品
			stock.PUT("/product", v1.UpdateStockProduct)                     // 更新产品
			stock.DELETE("/product/:id", v1.DeleteStockProduct)              // 删除产品

			// 产品明细相关
			stock.GET("/product-details", v1.GetStockProductDetails)               // 获取明细列表
			stock.GET("/product-details/summary", v1.GetStockProductDetailSummary) // 获取明细汇总
			stock.GET("/product-detail/:id", v1.GetStockProductDetail)             // 获取单个明细
			stock.POST("/product-detail", v1.AddStockProductDetail)                // 创建明细
			stock.PUT("/product-detail", v1.UpdateStockProductDetail)              // 更新明细
			stock.DELETE("/product-detail/:id", v1.DeleteStockProductDetail)       // 删除明细

			// 事务相关
			stock.POST("/transaction/batch-create", v1.BatchCreateProductWithDetails) // 批量创建产品及明细（事务）
			stock.POST("/transaction/transfer", v1.TransferStock)                     // 库存转移（事务）

			// Hooks 相关
			stock.POST("/hooks/product", v1.CreateProductWithHooks)    // 创建产品（演示 Hooks）
			stock.PUT("/hooks/product/:id", v1.UpdateProductWithHooks) // 更新产品（演示 Hooks）

			// 关联查询优化
			stock.GET("/optimize/products-with-details", v1.GetProductsWithDetailsOptimized) // 预加载关联数据
			stock.GET("/optimize/products-join", v1.GetProductsWithJoin)                     // 使用 Join 查询

			// 事务嵌套
			stock.POST("/nested-transaction/order", v1.CreateOrderWithNestedTransaction) // 创建订单（嵌套事务）

			// 协程并发处理
			stock.POST("/goroutine/batch-query", v1.BatchQueryProductsWithGoroutine) // 并发查询多个产品
			stock.POST("/goroutine/batch-update", v1.BatchUpdateStockWithGoroutine)  // 并发更新库存
			stock.POST("/goroutine/worker-pool", v1.ProcessWithWorkerPool)           // Worker Pool 模式
			stock.GET("/goroutine/pipeline", v1.ProcessWithPipeline)                 // Pipeline 模式
			stock.POST("/goroutine/fan-out-in", v1.ProcessWithFanOutFanIn)           // Fan-out/Fan-in 模式
		}

		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	return r
}
