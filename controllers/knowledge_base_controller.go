package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"yuqueppbackend/dao"
	"yuqueppbackend/models"
)

// CreateKnowledgeBase 创建知识库
func CreateKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id          int64  `json:"userid"`
		Name        string `json:"kb_name" binding:"required"`
		Description string `json:"kb_description" binding:"required"`
		IsPublic    bool   `json:"kb_is_public"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}

	if err := c.ShouldBindJSON(&contextData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	knowledgeBase := models.KnowledgeBase{
		Name:        contextData.Name,
		Description: contextData.Description,
		IsPublic:    contextData.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OwnerID:     contextData.Id,
	}
	// 使用 KBDAO 来创建知识库
	kbDAO := dao.NewKBDAO()
	if err := kbDAO.CreateKB(knowledgeBase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create knowledge base"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kb_id":          knowledgeBase.ID,
		"kb_name":        knowledgeBase.Name,
		"kb_description": knowledgeBase.Description,
		"kb_is_public":   knowledgeBase.IsPublic,
		"kb_created_at":  knowledgeBase.CreatedAt,
		"kb_updated_at":  knowledgeBase.UpdatedAt,
	})
}

// 获取用户创建的所有知识库
func GetKnowledgeBaseList(c *gin.Context) {
	kbDAO := dao.NewKBDAO()
	userId, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}
	if kbList, err := kbDAO.GetKBListByOwnerId(userId.(int64)); err == nil {
		c.JSON(http.StatusOK, gin.H{"knowledge_bases": kbList})
		log.Println(kbList)
		return
	} else {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
	}
	return
}

// GetKnowledgeBaseDetail 根据用户ID和知识库ID获取知识库详情
func GetKnowledgeBaseDetail(c *gin.Context) {
	var contextData struct {
		Id   int64 `json:"userid"`
		KBId int64 `json:"kb_id" binding:"required"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println("结构绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 使用 KBDAO 查找知识库
	kbDAO := dao.NewKBDAO()
	knowledgeBase, err := kbDAO.FindKB(contextData.Id, contextData.KBId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Knowledge base not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kb_id":          knowledgeBase.ID,
		"kb_name":        knowledgeBase.Name,
		"kb_description": knowledgeBase.Description,
		"kb_is_public":   knowledgeBase.IsPublic,
		"kb_created_at":  knowledgeBase.CreatedAt,
		"kb_updated_at":  knowledgeBase.UpdatedAt,
	})
}

// UpdateKnowledgeBase 更新知识库,可以更新的字段：Name,Description,IsPublic
func UpdateKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id          int64     `json:"userid"`
		KBId        int64     `json:"kb_id" binding:"required"`
		Name        string    `json:"kb_name" binding:"required"`
		Description string    `json:"kb_description" binding:"required"`
		IsPublic    bool      `json:"kb_is_public" binding:"required"`
		updated_at  time.Time `json:"kb_updated_at"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println("结构绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contextData.updated_at = time.Now()
	var updatedKB models.KnowledgeBase

	if err := c.ShouldBindJSON(&updatedKB); err != nil {
		log.Println("结构绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 KBDAO 查找并更新知识库
	kbDAO := dao.NewKBDAO()
	knowledgeBase, err := kbDAO.UpdateKB(contextData.Id, contextData.KBId, updatedKB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update knowledge base"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kb_id":          knowledgeBase.ID,
		"kb_name":        knowledgeBase.Name,
		"kb_description": knowledgeBase.Description,
		"kb_is_public":   knowledgeBase.IsPublic,
		"kb_created_at":  knowledgeBase.CreatedAt,
		"kb_updated_at":  knowledgeBase.UpdatedAt,
	})
}

// DeleteKnowledgeBase 删除知识库
func DeleteKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id   int64 `json:"userid"`
		KBId int64 `json:"kb_id" binding:"required"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.Id = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println("结构绑定失败")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 使用 KBDAO 查找知识库
	kbDAO := dao.NewKBDAO()
	knowledgeBase, err := kbDAO.FindKB(contextData.Id, contextData.KBId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Knowledge base not found"})
	}
	// 使用 KBDAO 删除知识库
	if err := kbDAO.DeleteKB(knowledgeBase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete knowledge base"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
