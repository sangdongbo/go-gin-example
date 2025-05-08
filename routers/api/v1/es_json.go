package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/es"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

func AddEsJsonData(c *gin.Context) {
	appG := app.Gin{C: c}
	index := "products"

	// 1. 创建索引并定义 mapping
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title":    map[string]interface{}{"type": "text"},
				"category": map[string]interface{}{"type": "keyword"},
				"price":    map[string]interface{}{"type": "float"},
				"brand": map[string]interface{}{
					"properties": map[string]interface{}{
						"name":   map[string]interface{}{"type": "keyword"},
						"origin": map[string]interface{}{"type": "keyword"},
					},
				},
				"tags": map[string]interface{}{"type": "keyword"},
				"comments": map[string]interface{}{ // 不用 nested，仅作为对象数组
					"type": "object",
				},
				"attributes": map[string]interface{}{ // 同上
					"type": "object",
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		appG.Response(http.StatusInternalServerError, -1, "编码 mapping 失败")
		return
	}

	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  &buf,
	}
	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil {
		appG.Response(http.StatusInternalServerError, -1, "创建索引失败")
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		appG.Response(http.StatusInternalServerError, -1, "索引已存在或创建失败")
		return
	}

	// 2. 插入 10 条模拟商品数据
	for i := 1; i <= 10; i++ {
		doc := map[string]interface{}{
			"title":    fmt.Sprintf("测试手机 %d", i),
			"category": "手机",
			"price":    rand.Float64()*3000 + 1000,
			"brand": map[string]interface{}{
				"name":   "测试品牌",
				"origin": "中国",
			},
			"tags": []string{"新品", "热销", "5G"},
			"comments": []map[string]interface{}{
				{
					"user":         fmt.Sprintf("user_%02d", i),
					"content":      "使用体验非常不错！",
					"score":        rand.Intn(2) + 4,
					"comment_date": time.Now().AddDate(0, 0, -rand.Intn(30)).Format("2006-01-02"),
				},
			},
			"attributes": []map[string]interface{}{
				{"name": "颜色", "value": "黑色"},
				{"name": "存储", "value": "128GB"},
			},
		}

		_, err := es.Index(index, "", doc)
		if err != nil {
			appG.Response(http.StatusInternalServerError, -1, fmt.Sprintf("插入文档 %d 失败: %v", i, err))
			return
		}
	}

	appG.Response(http.StatusOK, 200, "创建 products 索引并成功插入 10 条商品数据")
}

type Product struct {
	Title      string                   `json:"title"`
	Category   string                   `json:"category"`
	Price      float64                  `json:"price"`
	Brand      map[string]string        `json:"brand"`
	Tags       []string                 `json:"tags"`
	Comments   []map[string]interface{} `json:"comments"`
	Attributes []map[string]interface{} `json:"attributes"`
}

func SearchProductByKeyword(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.Query("q")

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "tags"},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		appG.Response(http.StatusBadRequest, -1, "查询构造失败")
		return
	}

	req := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  &buf,
	}
	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil || res.IsError() {
		appG.Response(http.StatusInternalServerError, -1, "搜索请求失败")
		return
	}
	defer res.Body.Close()

	// 原始 JSON -> 泛型 map
	var raw map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		appG.Response(http.StatusInternalServerError, -1, "结果解析失败")
		return
	}

	// 提取 hits
	hits := raw["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []Product

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)

		var p Product
		if err := json.Unmarshal(sourceBytes, &p); err == nil {
			results = append(results, p)
		}
	}

	appG.Response(http.StatusOK, 200, results)
}

type Product2 struct {
	ID         string                   `json:"id"`
	Title      string                   `json:"title"`
	Category   string                   `json:"category"`
	Price      float64                  `json:"price"`
	Brand      map[string]string        `json:"brand"`
	Tags       []string                 `json:"tags"`
	Comments   []map[string]interface{} `json:"comments"`
	Attributes []map[string]interface{} `json:"attributes"`
	Highlight  map[string][]string      `json:"highlight,omitempty"`
}

func SearchProductWithHighlight(c *gin.Context) {
	appG := app.Gin{C: c}
	keyword := c.DefaultQuery("q", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	from := (page - 1) * size

	// 构造带高亮的查询
	query := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"title", "tags"},
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
				"tags":  map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		appG.Response(http.StatusBadRequest, -1, "查询构造失败")
		return
	}

	req := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  &buf,
	}
	res, err := req.Do(context.Background(), es.ESClient)
	if err != nil || res.IsError() {
		appG.Response(http.StatusInternalServerError, -1, "搜索请求失败")
		return
	}
	defer res.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		appG.Response(http.StatusInternalServerError, -1, "结果解析失败")
		return
	}

	hits := raw["hits"].(map[string]interface{})["hits"].([]interface{})
	var results []Product2

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"]
		sourceBytes, _ := json.Marshal(source)

		var p Product2
		if err := json.Unmarshal(sourceBytes, &p); err == nil {
			p.ID = hitMap["_id"].(string)

			if hl, ok := hitMap["highlight"]; ok {
				p.Highlight = map[string][]string{}
				for field, val := range hl.(map[string]interface{}) {
					strs := []string{}
					for _, v := range val.([]interface{}) {
						strs = append(strs, v.(string))
					}
					p.Highlight[field] = strs
				}
			}

			results = append(results, p)
		}
	}

	appG.Response(http.StatusOK, 200, gin.H{
		"list":  results,
		"total": raw["hits"].(map[string]interface{})["total"],
		"page":  page,
		"size":  size,
	})
}

func SearchByBrandOrigin(c *gin.Context) {
	appG := app.Gin{C: c}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"brand.name": "测试品牌"}},
					{"match": map[string]interface{}{"brand.origin": "中国"}},
				},
			},
		},
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(query)

	res, err := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  &buf,
	}.Do(context.Background(), es.ESClient)
	if err != nil || res.IsError() {
		appG.Response(http.StatusInternalServerError, -1, "品牌查询失败")
		return
	}
	defer res.Body.Close()

	var result map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&result)

	appG.Response(http.StatusOK, 200, result)
}

func FilterProductByPrice(c *gin.Context) {
	appG := app.Gin{C: c}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": 2000,
					"lte": 4000,
				},
			},
		},
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(query)

	res, err := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  &buf,
	}.Do(context.Background(), es.ESClient)
	if err != nil || res.IsError() {
		appG.Response(http.StatusInternalServerError, -1, "区间查询失败")
		return
	}
	defer res.Body.Close()

	var result map[string]interface{}
	_ = json.NewDecoder(res.Body).Decode(&result)

	appG.Response(http.StatusOK, 200, result)
}
