package dao

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
	"yuqueppbackend/models"
	"yuqueppbackend/util"
)

// DocDao 处理与 Document 表相关的数据库操作
type DocDao struct {
	db *gorm.DB
}

// NewDocDao 创建一个新的 DocDao 实例
func NewDocDao(db *gorm.DB) *DocDao {
	return &DocDao{db: db}
}

// CreateDocument 创建文档
func (dao *DocDao) CreateDocument(doc *models.Document) error {
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	return dao.db.Create(doc).Error
}

// GetDocumentByID 根据 ID 获取文档
func (dao *DocDao) GetDocumentByID(id int64) (*models.Document, error) {
	var doc models.Document
	err := dao.db.Preload("KnowledgeBase").Preload("Children").First(&doc, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// GetDocumentsByKnowledgeBaseID 获取某知识库下的所有文档
func (dao *DocDao) GetDocumentsByKnowledgeBaseID(kbID int64) ([]models.Document, error) {
	var docs []models.Document
	err := dao.db.Where("knowledge_base_id = ?", kbID).Find(&docs).Error
	return docs, err
}

// UpdateDocument 更新文档
func (dao *DocDao) UpdateDocument(doc *models.Document) error {
	doc.UpdatedAt = time.Now()
	return dao.db.Save(doc).Error
}

// DeleteDocumentByID 根据 ID 删除文档
func (dao *DocDao) DeleteDocumentByID(id int64) error {
	return dao.db.Delete(&models.Document{}, id).Error
}

// IncrementViewCount 增加文档的浏览次数
func (dao *DocDao) IncrementViewCount(id uint) error {
	return dao.db.Model(&models.Document{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// AddTag 为文档添加标签
func (dao *DocDao) AddTag(id int64, tag string) error {
	var doc models.Document
	err := dao.db.First(&doc, id).Error
	if err != nil {
		return err
	}

	// 追加标签
	if doc.Tags == "" {
		doc.Tags = tag
	} else {
		doc.Tags += "," + tag
	}

	doc.UpdatedAt = time.Now()
	return dao.db.Save(&doc).Error
}

// 向 redis中写入文档浏览记录
//func (dao *DocDao) UpdateRecentDocumentViewInRedis(userId, docId string) error {
//	// 写入 Redis 最近浏览记录
//	key := "user_recent_docs:" + userId
//	timestamp := float64(time.Now().Unix())
//	util.GetRedisClient().ZAdd(context.Background(), key, &redis.Z{
//		Score:  timestamp,
//		Member: docId,
//	})
//	// 限制最近浏览记录的条数
//	maxRecentCount := 50
//	util.GetRedisClient().ZRemRangeByRank(context.Background(), key, 0, int64(-maxRecentCount-1))
//	return nil
//}

// UpdateRecentDocumentViewInRedis 更新最近浏览记录到 Redis
func (dao *DocDao) UpdateRecentDocumentViewInRedis(document models.Document, kbName, userId string) error {
	key := "user_recent_docs:" + userId
	timestamp := float64(time.Now().Unix())

	// 构造 JSON 数据作为 Redis ZSet 的 Member
	member := map[string]interface{}{
		"doc_id":    document.ID,
		"doc_title": document.Title,
		"kb_id":     document.KnowledgeBaseID,
		"kb_name":   kbName,
	}
	memberJSON, err := json.Marshal(member)
	if err != nil {
		return err
	}

	// 向 Redis ZSet 添加数据
	util.GetRedisClient().ZAdd(context.Background(), key, &redis.Z{
		Score:  timestamp,
		Member: memberJSON,
	})

	// 限制最多保留 50 条最近记录
	maxRecentCount := 50
	util.GetRedisClient().ZRemRangeByRank(context.Background(), key, 0, int64(-maxRecentCount-1))

	return nil
}

// GetRecentDocumentViewsWithScoresFromRedis 获取用户最近浏览的文档记录（带分数）
func (dao *DocDao) GetRecentDocumentViewsWithScoresFromRedis(userId string, start, end int64) ([]map[string]interface{}, error) {
	key := "user_recent_docs:" + userId

	// 获取带分数的数据
	values, err := util.GetRedisClient().ZRevRangeWithScores(context.Background(), key, start, end).Result()
	if err != nil {
		return nil, err
	}

	// 解析数据
	var recentDocuments []map[string]interface{}
	for _, value := range values {
		var document map[string]interface{}
		if err := json.Unmarshal([]byte(value.Member.(string)), &document); err != nil {
			// 如果解析失败，记录错误但继续处理其他记录
			log.Printf("Failed to parse recent document entry: %v", err)
			continue
		}
		// 附加分数（时间戳）
		document["timestamp"] = value.Score
		recentDocuments = append(recentDocuments, document)
	}

	return recentDocuments, nil
}
