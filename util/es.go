package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strconv"
	"yuqueppbackend/config"
	"yuqueppbackend/models"
)

var esClient *elasticsearch.Client

func InitElasticSearchClient() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.GetElasticSearchAddress(),
		},
	}
	var err error
	esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 测试连接
	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

// 将文档数据插入到ES

func InsertDocToES(document models.Document, content string) error {
	docIdStr := strconv.FormatInt(int64(document.ID), 10)
	doc := map[string]interface{}{
		"id":      docIdStr,
		"title":   document.Title,
		"content": content,
	}

	// 将文档转为json
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}

	// 索引文档
	res, err := esClient.Index(
		"document",
		&buf,
		esClient.Index.WithDocumentID(docIdStr), // 自定义 _id
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	fmt.Println(res)
	return nil
}

func UpdateDocToES(documentID int64, title string, content string) error {
	// 创建需要更新的字段数据
	doc := map[string]interface{}{
		"doc": map[string]interface{}{
			"title":   title,
			"content": content,
		},
	}

	// 将更新数据转为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding update data: %s", err)
	}

	// 执行更新操作
	res, err := esClient.Update(
		"document",                        // 索引名称
		strconv.FormatInt(documentID, 10), // 文档 ID（字符串）
		&buf,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}

func DeleteDocFromES(documentID int64) error {
	// 执行删除操作
	res, err := esClient.Delete(
		"document",                        // 索引名称
		strconv.FormatInt(documentID, 10), // 文档 ID
	)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
	return nil
}

func GetDocIDByTitle(title string) ([]string, error) {
	// 构造查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": title,
			},
		},
	}

	// 将查询条件转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索
	res, err := esClient.Search(
		esClient.Search.WithIndex("document"),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 提取文档 ID
	var ids []string
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		for _, hit := range hits["hits"].([]interface{}) {
			if doc, ok := hit.(map[string]interface{}); ok {
				ids = append(ids, doc["_id"].(string))
			}
		}
	}

	return ids, nil
}

func GetDocIDByContent(content string) ([]string, error) {
	// 构造查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": content,
			},
		},
	}

	// 将查询条件转换为 JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// 执行搜索
	res, err := esClient.Search(
		esClient.Search.WithIndex("document"),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 解析结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 提取文档 ID
	var ids []string
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		for _, hit := range hits["hits"].([]interface{}) {
			if doc, ok := hit.(map[string]interface{}); ok {
				ids = append(ids, doc["_id"].(string))
			}
		}
	}

	return ids, nil
}
