package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yuqueppbackend/dao"
	"yuqueppbackend/models"
)

type DocumentController struct {
	docDao *dao.DocDao
}

// NewDocumentController 创建新的 DocumentController
func NewDocumentController(docDao *dao.DocDao) *DocumentController {
	return &DocumentController{docDao: docDao}
}

// CreateDocumentHandler 创建文档
func (dc *DocumentController) CreateDocumentHandler(c *gin.Context) {

	var contextData struct {
		UserId int64  `json:"userid"`
		DocId  int64  `json:"doc_id" binding:"required"`
		Title  string `json:"doc_title"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.UserId = id.(int64)
	}

	if err := c.ShouldBindJSON(&contextData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	if contextData.Title == "" {
		contextData.Title = "无标题"
	}
	var doc models.Document = models.Document{
		ID:      contextData.DocId,
		Title:   contextData.Title,
		OwnerId: contextData.UserId,
	}

	if err := dc.docDao.CreateDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"doc_id":      doc.ID,
		"kb_id":       doc.KnowledgeBaseID,
		"doc_title":   doc.Title,
		"doc_content": doc.Content,
	})
}

// GetDocumentByIDHandler 获取文档详情
func (dc *DocumentController) GetDocumentByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := dc.docDao.GetDocumentByID(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
		return
	}
	if doc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"doc_id":      doc.ID,
		"kb_id":       doc.KnowledgeBaseID,
		"doc_title":   doc.Title,
		"doc_content": doc.Content,
	})
}

// GetDocumentsByKnowledgeBaseIDHandler 获取某知识库的所有文档
func (dc *DocumentController) GetDocumentsByKnowledgeBaseIDHandler(c *gin.Context) {
	kbIDParam := c.Param("kb_id")
	kbID, err := strconv.Atoi(kbIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid knowledge base ID"})
		return
	}

	docs, err := dc.docDao.GetDocumentsByKnowledgeBaseID(int64(kbID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}
	docList := make(map[string]interface{}, 0)
	for _, doc := range docs {
		tmpMap := make(map[string]interface{})
		tmpMap["kb_id"] = doc.KnowledgeBaseID
		tmpMap["doc_id"] = doc.ID
		tmpMap["doc_title"] = doc.Title
	}
	c.JSON(http.StatusOK, docList)
}

// UpdateDocumentHandler 更新文档
func (dc *DocumentController) UpdateDocumentHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var doc models.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	doc.ID = int64(id)
	if err := dc.docDao.UpdateDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// DeleteDocumentByIDHandler 删除文档
func (dc *DocumentController) DeleteDocumentByIDHandler(c *gin.Context) {
	// 从请求参数获取文档ID
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// 获取当前登录用户的ID（假设JWT认证中间件已经将用户信息存储在上下文中）
	userID, exists := c.Get("userid") // 从上下文中获取当前用户的ID
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 查询文档的所有者
	document, err := dc.docDao.GetDocumentByID(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
		return
	}

	// 判断当前用户是否为文档的所有者
	if document.OwnerId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this document"})
		return
	}

	// 执行删除操作
	if err := dc.docDao.DeleteDocumentByID(int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// IncrementViewCountHandler 增加浏览次数
func (dc *DocumentController) IncrementViewCountHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	if err := dc.docDao.IncrementViewCount(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment view count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View count incremented"})
}
