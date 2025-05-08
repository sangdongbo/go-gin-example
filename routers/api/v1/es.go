package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/es"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

func AddEsData(c *gin.Context) {
	appG := app.Gin{C: c}
	index := "articles"

	// 1. 创建复杂索引（带字段类型定义）
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
				},
				"content": map[string]interface{}{
					"type": "text",
				},
				"author": map[string]interface{}{
					"type": "keyword",
				},
				"tags": map[string]interface{}{
					"type": "keyword",
				},
				"published_at": map[string]interface{}{
					"type": "date",
				},
				"views": map[string]interface{}{
					"type": "integer",
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		appG.Response(http.StatusInternalServerError, -1, "failed to encode mapping")
		return
	}

	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  &buf,
	}
	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil {
		appG.Response(http.StatusInternalServerError, -1, "failed to create index")
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		appG.Response(http.StatusInternalServerError, -1, "index already exists or error occurred")
		return
	}

	// 2. 插入 100 行数据
	authors := []string{"Alice", "Bob", "Charlie", "Diana"}
	tags := [][]string{
		{"tech", "go"},
		{"life", "travel"},
		{"food", "recipe"},
		{"news", "world"},
	}

	for i := 1; i <= 100; i++ {
		doc := map[string]interface{}{
			"title":        fmt.Sprintf("第 %d 篇文章", i),
			"content":      fmt.Sprintf("这是第 %d 篇文章的内容，非常精彩……", i),
			"author":       authors[i%len(authors)],
			"tags":         tags[i%len(tags)],
			"published_at": time.Now().AddDate(0, 0, -rand.Intn(365)), // 随机过去 1 年内
			"views":        rand.Intn(10000),
		}

		_, err := es.Index(index, "", doc)
		if err != nil {
			appG.Response(http.StatusInternalServerError, -1, fmt.Sprintf("failed to insert doc %d: %v", i, err))
			return
		}
	}

	appG.Response(http.StatusOK, 200, "索引创建并插入 100 条数据成功")
}

func GetEsData(c *gin.Context) {
	appG := app.Gin{C: c}

	// 获取查询参数
	keyword := c.DefaultQuery("keyword", "")
	author := c.Query("author")
	tag := c.Query("tag")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := 10
	from := (page - 1) * pageSize

	// 构造查询
	must := []map[string]interface{}{}
	filter := []map[string]interface{}{}

	if keyword != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "author", "content"},
			},
		})
	}

	if author != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"author": author,
			},
		})
	}

	if tag != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"tags": tag,
			},
		})
	}

	// 时间范围：过去 180 天
	filter = append(filter, map[string]interface{}{
		"range": map[string]interface{}{
			"published_at": map[string]interface{}{
				"gte": time.Now().AddDate(0, 0, -180).Format(time.RFC3339),
			},
		},
	})

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"sort": []map[string]interface{}{
			{"views": map[string]string{"order": "desc"}},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
		},
	}

	results, err := es.Search("articles", query, from, pageSize)
	if err != nil {
		appG.Response(http.StatusInternalServerError, -1, err.Error())
		return
	}

	appG.Response(http.StatusOK, 0, results)
}
