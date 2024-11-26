package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"yuqueppbackend/config"
	"yuqueppbackend/dao"
	"yuqueppbackend/models"
)

type DocumentController struct {
	docDao *dao.DocDao
}

func getDocumentContentById(docId string) (string, error) {
	content, err := os.ReadFile(config.GetDocumentStoragePath() + "/" + docId + ".txt")
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(content), nil
}
func deleteDocumentFile(docId string) error {
	err := os.Remove(config.GetDocumentStoragePath() + "/" + docId + ".txt")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// NewDocumentController 创建新的 DocumentController
func NewDocumentController(docDao *dao.DocDao) *DocumentController {
	return &DocumentController{docDao: docDao}
}

// CreateDocumentHandler 创建文档
func (dc *DocumentController) CreateDocumentHandler(c *gin.Context) {

	var contextData struct {
		UserId int64  `json:"userid"`
		KbId   string `json:"kb_id" binding:"required"`
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
	kbId64, err := strconv.ParseInt(contextData.KbId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	var doc models.Document = models.Document{
		KnowledgeBaseID: kbId64,
		Title:           contextData.Title,
		OwnerId:         contextData.UserId,
	}

	if err := dc.docDao.CreateDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}
	str_doc_id := strconv.FormatInt(doc.ID, 10)
	doc_file, err := os.Create(config.GetDocumentStoragePath() + "/" + str_doc_id + ".txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	// 确保文件内容为空
	err = doc_file.Truncate(0)
	if _, err := doc_file.Write([]byte("# " + doc.Title)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法清空文件内容"})
		return
	}
	// 关闭文件
	defer doc_file.Close()

	c.JSON(http.StatusOK, gin.H{
		"doc_id":      str_doc_id,
		"kb_id":       doc.KnowledgeBaseID,
		"doc_title":   doc.Title,
		"doc_content": doc.Content,
	})
}

// GetDocumentByIDHandler 获取文档详情
func (dc *DocumentController) GetDocumentByIDHandler(c *gin.Context) {
	idParam := c.Param("doc_id")
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
		c.JSON(http.StatusNotFound, gin.H{"error": "当前文档不见了，快去新建吧"})
		return
	}
	strDocId := strconv.FormatInt(doc.ID, 10)
	strKbId := strconv.FormatInt(doc.KnowledgeBaseID, 10)
	docContent, err := getDocumentContentById(strDocId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，请稍后再试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"doc_id":      strDocId,
		"kb_id":       strKbId,
		"doc_title":   doc.Title,
		"doc_content": docContent,
	})
	return
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
	var docList []map[string]interface{}
	for _, doc := range docs {
		tmpMap := make(map[string]interface{})
		tmpMap["kb_id"] = strconv.FormatInt(doc.KnowledgeBaseID, 10)
		tmpMap["doc_id"] = strconv.FormatInt(doc.ID, 10)
		tmpMap["doc_title"] = doc.Title
		docList = append(docList, tmpMap)
	}
	c.JSON(http.StatusOK, gin.H{"doc_list": docList})
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
	idParam := c.Param("doc_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "错误的文档ID"})
		return
	}

	userID, exists := c.Get("userid") // 从上下文中获取当前用户的ID
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户没有权限"})
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
	deleteDocumentFile(idParam)

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

func (dc *DocumentController) GetDocumentListByKbId(c *gin.Context) {
	// 定义请求上下文结构
	var contextData struct {
		UserID int64  `json:"userid"`
		KBId   string `json:"kb_id" binding:"required"`
	}

	// 获取用户 ID（从中间件中设置的上下文）
	if id, exists := c.Get("userid"); exists {
		contextData.UserID = id.(int64)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未授权"})
		return
	}

	// 绑定 JSON 请求参数
	if err := c.ShouldBindJSON(&contextData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 解析 kb_id 参数
	kbID, err := strconv.ParseInt(contextData.KBId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的知识库 ID"})
		return
	}

	// 查询文档列表
	var documents []models.Document
	docList, err := dc.docDao.GetDocumentsByKnowledgeBaseID(kbID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，获取文档失败"})
		return
	}
	var docListData []map[string]interface{}
	for _, doc := range docList {
		tmpMap := make(map[string]interface{})
		tmpMap["doc_id"] = strconv.FormatInt(doc.ID, 10)
		tmpMap["doc_created_at"] = doc.CreatedAt
		tmpMap["doc_updated_at"] = doc.UpdatedAt
		docListData = append(docListData, tmpMap)
	}

	c.JSON(http.StatusOK, gin.H{"doc_list": docListData})
	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}
