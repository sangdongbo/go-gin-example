package v1

import (
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/es"
	"github.com/EDDYCJY/go-gin-example/service/rabbitmq_service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func AddEmails(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form rabbitmq_service.AddBaseEmailForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	email := rabbitmq_service.ConvertAddFormToUEmail(form)
	id, err := email.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ORDER_FAIL, nil)
		return
	}

	// 创建 ES 索引
	esDoc := map[string]interface{}{
		"user_id": id,
		"subject": form.Subject,
		"body":    form.Body,
		"status":  form.Status,
	}

	docID := strconv.Itoa(id)

	if _, err := es.Index("email_index", docID, esDoc); err != nil {
		// 索引失败不阻断主流程，但记录日志
		log.Printf("failed to index email to ES: %v", err)
	}

	//添加信息的 es 中，创建索引
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func UpdateEmail(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form rabbitmq_service.UpdateBaseEmailForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	email := rabbitmq_service.ConvertEditFormToUEmail(form)

	// 先更新数据库
	err := email.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ORDER_FAIL, nil)
		return
	}

	// 然后更新 Elasticsearch 的文档
	esDocID := strconv.Itoa(email.ID) // 用 ID 作为 Elasticsearch 的文档 ID
	index := "email_index"            // 确保这是你在 ES 中用的索引名

	_, err = es.Index(index, esDocID, email)
	if err != nil {
		// ES 更新失败，不影响主流程，但可以记录日志或设置告警
		log.Printf("failed to update email in ES, ID=%d: %v", email.ID, err)
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// Search 执行带分页的查询
func GetEmails(c *gin.Context) {
	appG := app.Gin{C: c}

	// 获取分页参数
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	from := (pageNum - 1) * pageSize

	// 其他查询参数
	userIDStr := c.Query("user_id")
	status := c.Query("status")
	keyword := c.Query("q")

	// 构造 Elasticsearch 查询
	boolQuery := map[string]interface{}{
		"must": []interface{}{},
	}

	// 添加条件
	if userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
				"term": map[string]interface{}{"user_id": userID},
			})
		}
	}
	if status != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{"status": status},
		})
	}
	if keyword != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"subject", "body"},
			},
		})
	}

	// 完整 query 对象
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": boolQuery,
		},
	}

	// 调用封装的 Search 方法
	results, err := es.Search("email_index", query, from, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, results)
}
