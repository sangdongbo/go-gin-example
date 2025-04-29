package v1

import (
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/order_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "OrderSn"
// @Param state query int false "UserId"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/orders [get]
func GetOrders(c *gin.Context) {
	appG := app.Gin{C: c}
	//appG := app.Gin{C: c}
	orderSn := c.Query("order_sn")

	orderService := order_service.Order{
		OrderSn:  orderSn,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	orders, err := orderService.GetAll()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ORDER_FAIL, nil)
		return
	}

	count, err := orderService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ORDER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": orders,
		"total": count,
	})
}

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

	orderService := order_service.ConvertAddOrderFormToOrder(form)

	err := orderService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
