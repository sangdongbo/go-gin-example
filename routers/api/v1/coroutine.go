package v1

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

// TestCoroutineRequest 测试协程请求结构
type TestCoroutineRequest struct {
	IDs      []int `json:"ids" binding:"required"`
	TestMode *bool `json:"test_mode"` // 默认为 true
}

// TestCoroutineResponse 测试协程响应结构
type TestCoroutineResponse struct {
	ID   int                    `json:"id"`
	Data map[string]interface{} `json:"data"`
}

// @Summary 测试协程处理
// @Description 模拟 PHP Swoole 协程处理多个 ID
// @Produce  json
// @Param body body TestCoroutineRequest true "请求参数"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/coroutine/test [post]
func TestCoroutine(c *gin.Context) {
	appG := app.Gin{C: c}

	var req TestCoroutineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 日志文件路径
	logPath := "runtime/logs/debug.log"

	if req.TestMode != nil && *req.TestMode {
		writeLog(logPath, "使用测试模式（不访问数据库）")

		// 创建一个 channel 用于收集结果
		resultChan := make(chan TestCoroutineResponse, len(req.IDs))

		// 使用 WaitGroup 等待所有协程完成
		var wg sync.WaitGroup

		// 启动多个 goroutine 处理每个 ID
		for _, id := range req.IDs {
			wg.Add(1)
			func(id int) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						writeLog(logPath, fmt.Sprintf("测试协程错误 ID %d: %v", id, r))
						resultChan <- TestCoroutineResponse{
							ID:   id,
							Data: map[string]interface{}{},
						}
					}
				}()

				writeLog(logPath, fmt.Sprintf("测试协程 ID: %d", id))

				// 模拟一些工作
				time.Sleep(3 * time.Second)

				// 将结果推送到 channel
				resultChan <- TestCoroutineResponse{
					ID: id,
					Data: map[string]interface{}{
						"test": true,
						"id":   id,
					},
				}
				writeLog(logPath, fmt.Sprintf("测试协程完成 ID: %d", id))
			}(id)
		}

		writeLog(logPath, "等待测试协程完成...")

		// 等待所有协程完成
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// 收集结果
		data := make(map[int]map[string]interface{})
		timeout := time.After(5 * time.Second)

		for {
			select {
			case result, ok := <-resultChan:
				if !ok {
					// channel 已关闭，所有结果都已收集
					writeLog(logPath, "测试协程全部完成")
					appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
						"message": "协程处理完成",
						"data":    data,
						"count":   len(data),
					})
					return
				}
				data[result.ID] = result.Data
				writeLog(logPath, fmt.Sprintf("收到测试结果 ID: %d", result.ID))
			case <-timeout:
				writeLog(logPath, "测试协程超时")
				appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
					"message": "部分协程超时",
					"data":    data,
					"count":   len(data),
				})
				return
			}
		}
	}

	// 非测试模式
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"message": "请启用测试模式",
	})
}

// writeLog 写入日志到文件
func writeLog(path, message string) {
	logMessage := fmt.Sprintf("%s - %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开日志文件:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logMessage); err != nil {
		fmt.Println("写入日志失败:", err)
	}
}
