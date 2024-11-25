package dao

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"yuqueppbackend/models"
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
