package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

var ESClient *elasticsearch.Client

// Setup 初始化 Elasticsearch 客户端
func Setup() error {
	cfg := elasticsearch.Config{
		Addresses: []string{setting.ElasticSearchSetting.Hosts},
		Username:  setting.ElasticSearchSetting.Username,
		Password:  setting.ElasticSearchSetting.Password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("failed to get Elasticsearch info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error connecting to Elasticsearch: %s", string(body))
	}

	ESClient = client
	return nil
}

// Index 索引文档，如果 id 为空则自动生成
func Index(index string, id string, body interface{}) (string, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		Body:       bytes.NewReader(data),
		DocumentID: id, // 留空则自动生成
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return "", fmt.Errorf("index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("error indexing document: %s", string(body))
	}

	// 获取生成的 ID（如果自动生成）
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode index response: %w", err)
	}

	docID, ok := result["_id"].(string)
	if !ok {
		return "", errors.New("missing _id in response")
	}

	return docID, nil
}

// Get 获取文档
func Get(index, id string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, nil
	}

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error getting document ID=%s: %s", id, string(body))
	}

	var doc map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode get response: %w", err)
	}

	source, ok := doc["_source"].(map[string]interface{})
	if !ok {
		return nil, errors.New("document missing _source")
	}

	return source, nil
}

// Delete 删除文档
func Delete(index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error deleting document ID=%s: %s", id, string(body))
	}

	return nil
}

// Search 执行带分页的查询
func Search(index string, query map[string]interface{}, from, size int) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := ESClient.Search(
		ESClient.Search.WithContext(context.Background()),
		ESClient.Search.WithIndex(index),
		ESClient.Search.WithBody(&buf),
		ESClient.Search.WithFrom(from),
		ESClient.Search.WithSize(size),
		ESClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error searching index=%s: %s", index, string(body))
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	hitsData, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		log.Printf("unexpected structure in search result: %+v", r["hits"])
		return nil, errors.New("invalid hits structure")
	}

	results := make([]map[string]interface{}, len(hitsData))
	for i, hit := range hitsData {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		results[i] = source
	}
	return results, nil
}

// UpdateReplica 设置索引副本数
func UpdateReplica(index string, numReplicas int) error {
	settings := map[string]interface{}{
		"index": map[string]interface{}{
			"number_of_replicas": numReplicas,
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(settings); err != nil {
		return fmt.Errorf("failed to encode settings: %w", err)
	}

	req := esapi.IndicesPutSettingsRequest{
		Index: []string{index},
		Body:  &buf,
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error updating settings: %s", string(body))
	}
	return nil
}

// Update 更新文档字段
func Update(index, id string, doc map[string]interface{}) error {
	updateBody := map[string]interface{}{
		"doc": doc,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(updateBody); err != nil {
		return fmt.Errorf("failed to encode update doc: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       &buf,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return fmt.Errorf("update request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error updating document ID=%s: %s", id, string(body))
	}

	return nil
}
