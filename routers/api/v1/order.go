package v1

import (
	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
)

// @Summary 获取订单列表
// @Tags 订单管理
// @Produce json
// @Param order_sn query string false "订单号"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/orders [get]
func GetOrders(c *gin.Context) {
	appG := app.Gin{C: c}

	query := order_service.OrderQuery{
		OrderSn:  c.Query("order_sn"),
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	orders, err := query.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ORDER_FAIL, nil)
		return
	}

	count, err := query.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ORDER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": orders,
		"total": count,
	})
}

// @Summary 创建订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param order body order_service.AddOrderForm true "订单信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order [post]
func AddOrder(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form order_service.AddOrderForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	order := order_service.ConvertAddFormToOrder(form)
	if err := order.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ORDER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 更新订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param order body order_service.EditOrderForm true "订单信息"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 404 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/order [put]
func UpdateOrder(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form order_service.EditOrderForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	order := order_service.ConvertEditFormToOrder(form)

	// 验证订单是否存在
	if _, err := order_service.GetOne(form.Id); err != nil {
		appG.Response(http.StatusNotFound, e.ERROR_EDIT_ORDER_FAIL, nil)
		return
	}

	if err := order.Edit(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ORDER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func DeleteOrder(c *gin.Context) {
	appG := app.Gin{C: c}

	deleteOrderService := order_service.OrderDelete{
		OrderSn: c.Param("order_sn"),
	}

	err := deleteOrderService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ORDER_FAIL, nil)
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除成功")
}
