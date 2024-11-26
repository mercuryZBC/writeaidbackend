package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
	"yuqueppbackend/dao"
	"yuqueppbackend/models"
)

// CreateKnowledgeBase 创建知识库
func CreateKnowledgeBase(c *gin.Context) {
	var contextData struct {
		Id          int64  `json:"userid"`
		Name        string `json:"kb_name" binding:"required"`
		Description string `json:"kb_description"`
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
		"kb_id":          strconv.FormatInt(contextData.Id, 10),
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

		var kbListData []map[string]interface{}
		for _, kb := range kbList {
			tmpMap := make(map[string]interface{})
			tmpMap["kb_id"] = strconv.FormatInt(kb.ID, 10)
			tmpMap["kb_name"] = kb.Name
			tmpMap["kb_description"] = kb.Description
			tmpMap["kb_is_public"] = kb.IsPublic
			tmpMap["kb_created_at"] = kb.CreatedAt
			tmpMap["kb_updated_at"] = kb.UpdatedAt
			kbListData = append(kbListData, tmpMap)
		}

		c.JSON(http.StatusOK, gin.H{"knowledge_bases": kbListData})
		return
	} else {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
	}
	return
}

// GetKnowledgeBaseDetail 根据用户ID和知识库ID获取知识库详情
func GetKnowledgeBaseDetail(c *gin.Context) {
	kbId := c.Param("kb_id")
	id, exists := c.Get("userid")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "系统错误请稍后再试"})
		return
	}
	kbId64, _ := strconv.ParseInt(kbId, 10, 64)
	userId64, _ := id.(int64)

	// 使用 KBDAO 查找知识库
	kbDAO := dao.NewKBDAO()
	knowledgeBase, err := kbDAO.FindKB(userId64, kbId64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Knowledge base not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kb_id":          strconv.FormatInt(kbId64, 10),
		"kb_owner_id":    strconv.FormatInt(knowledgeBase.OwnerID, 10),
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
		Id   int64  `json:"userid"`
		KBId string `json:"kb_id" binding:"required"`
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
	// 使用 strconv.ParseInt 将字符串转换为 int64
	kbId64, err := strconv.ParseInt(contextData.KBId, 10, 64) // 10 是十进制，64 表示返回 int64 类型
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误"})
		return
	}
	knowledgeBase, err := kbDAO.FindKB(contextData.Id, kbId64)
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
