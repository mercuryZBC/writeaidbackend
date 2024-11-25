package models

import "time"

// KnowledgeBase 模型，作为主表
type KnowledgeBase struct {
	ID          uint      `json:"kb_id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"kb_name" binding:"required"`
	Description string    `json:"kb_description"`
	IsPublic    bool      `json:"kb_is_public"` // 是否公开
	OwnerID     uint      `json:"kb_owner_id"`  // 所有者
	CreatedAt   time.Time `json:"kb_created_at"`
	UpdatedAt   time.Time `json:"kb_updated_at"`

	User User `json:"user" gorm:"foreignKey:OwnerID;references:ID"`
	// 一对多关系
	Documents []Document `json:"directories" gorm:"foreignKey:KnowledgeBaseID;constraint:OnDelete:CASCADE;"`
}

// Document 模型，表示知识库中的文档
type Document struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title           string    `json:"title" binding:"required"` // 文档标题
	Content         string    `json:"content"`                  // 文档内容
	KnowledgeBaseID uint      `json:"knowledge_base_id"`        // 外键，所属知识库
	ParentID        *uint     `json:"parent_id"`                // 自引用外键，父文档 ID
	Status          string    `json:"status"`                   // 文档状态（如草稿、发布等）
	Tags            string    `json:"tags"`                     // 标签，逗号分隔
	ViewCount       uint      `json:"view_count"`               // 浏览次数
	CommentCount    uint      `json:"comment_count"`            // 评论数量
	CreatedAt       time.Time `json:"created_at"`               // 创建时间
	UpdatedAt       time.Time `json:"updated_at"`               // 更新时间
	Type            string    `json:"type"`                     // 文档类型（如文章、教程、参考等）

	// 关联的知识库
	KnowledgeBase KnowledgeBase `json:"knowledge_base" gorm:"foreignKey:KnowledgeBaseID;references:ID"`

	// 子文档，适用于文档的层级结构
	Children []Document `json:"children" gorm:"foreignKey:ParentID"`
}
