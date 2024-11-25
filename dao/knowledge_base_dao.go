package dao

import (
	"fmt"
	"gorm.io/gorm"
	"yuqueppbackend/db"
	"yuqueppbackend/models"
)

type KBDAO struct {
	DB *gorm.DB
}

func NewKBDAO() *KBDAO {
	return &KBDAO{DB: db.GetDB()}
}

func (dao *KBDAO) CreateKB(kb models.KnowledgeBase) error {
	return dao.DB.Create(&kb).Error
}

func (dao *KBDAO) GetKBListByOwnerId(ownerId uint) ([]models.KnowledgeBase, error) {
	var knowledgeBases []models.KnowledgeBase
	if err := dao.DB.Where("owner_id", ownerId).Find(&knowledgeBases).Error; err != nil {
		return nil, err
	}
	return knowledgeBases, nil
}

func (dao *KBDAO) DeleteKB(kb models.KnowledgeBase) error {
	return dao.DB.Delete(&kb).Error
}

// FindKB 查找知识库，使用 ownerId 和 id 作为查询条件
func (dao *KBDAO) FindKB(ownerId, id uint) (kb models.KnowledgeBase, err error) {
	// 查询知识库，使用 ownerId 和 id 作为条件
	err = dao.DB.Where("owner_id = ? AND id = ?", ownerId, id).First(&kb).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到记录，返回一个适当的错误
			return kb, fmt.Errorf("knowledge base not found for ownerId: %d, id: %d", ownerId, id)
		}
		// 其他数据库查询错误
		return kb, fmt.Errorf("failed to find knowledge base: %v", err)
	}

	// 返回查询到的知识库对象
	return kb, nil
}

// UpdateKB 更新知识库信息，使用 ownerId 和 id 作为查询条件
func (dao *KBDAO) UpdateKB(ownerId, id uint, updatedKB models.KnowledgeBase) (models.KnowledgeBase, error) {
	var kb models.KnowledgeBase

	// 查找指定 ownerId 和 id 的知识库
	err := dao.DB.Where("owner_id = ? AND id = ?", ownerId, id).First(&kb).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到记录，返回一个适当的错误
			return kb, fmt.Errorf("knowledge base not found for ownerId: %d, id: %d", ownerId, id)
		}
		// 其他数据库查询错误
		return kb, fmt.Errorf("failed to find knowledge base: %v", err)
	}

	// 更新字段，只更新传入的字段
	err = dao.DB.Model(&kb).Updates(updatedKB).Error
	if err != nil {
		// 更新失败，返回错误
		return kb, fmt.Errorf("failed to update knowledge base: %v", err)
	}

	// 返回更新后的知识库对象
	return kb, nil
}
