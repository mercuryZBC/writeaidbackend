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

type CommentController struct {
	commentDao *dao.CommentDAO
}

func NewCommentController(commentDao *dao.CommentDAO) *CommentController {
	return &CommentController{commentDao: commentDao}
}

func (cc *CommentController) ReplyDocumentComment(c *gin.Context) {
	var contextData struct {
		UserId         int64  `json:"userid"`
		DocId          string `json:"doc_id" binding:"required"`
		RootId         string `json:"root_id" binding:"required"`
		ParentId       string `json:"parent_id"`
		CommentContent string `json:"comment_content"`
	}

	if id, exists := c.Get("userid"); exists {
		contextData.UserId = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，评论创建失败"})
		return
	}
	docId, err := strconv.ParseInt(contextData.DocId, 10, 64)
	if err != nil {
		log.Println("数据格式转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，评论创建失败"})
		return
	}
	var rootId *int64 = nil
	if contextData.ParentId != "" {
		parsedId, err := strconv.ParseInt(contextData.RootId, 10, 64)
		if err == nil {
			rootId = &parsedId // 将 parsedId 的地址赋给 parentId
		}
	}
	var parentId *int64 = nil
	if contextData.ParentId != "" {
		parsedId, err := strconv.ParseInt(contextData.ParentId, 10, 64)
		if err == nil {
			parentId = &parsedId // 将 parsedId 的地址赋给 parentId
		}
	}

	dc := models.DocumentComment{
		DocumentID:   docId,
		ParentID:     parentId,
		RootID:       rootId,
		UserID:       contextData.UserId,
		Content:      contextData.CommentContent,
		Status:       "已发布",
		LikeCount:    0,
		DislikeCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
		CreatedAtBy:  "",
		EditedAtBy:   "",
	}

	err = cc.commentDao.CreateComment(&dc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "系统错误，评论回复失败"})
	}
	err = cc.commentDao.InsertReplyCommentToRedis(dc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "系统错误，评论回复失败"})
	}
	c.JSON(http.StatusOK, gin.H{})

}

func (cc *CommentController) CreateDocumentComment(c *gin.Context) {
	var contextData struct {
		UserId         int64  `json:"userid"`
		DocId          string `json:"doc_id" binding:"required"`
		CommentContent string `json:"comment_content"`
	}
	if id, exists := c.Get("userid"); exists {
		contextData.UserId = id.(int64)
	}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，评论创建失败"})
		return
	}
	docId, err := strconv.ParseInt(contextData.DocId, 10, 64)
	if err != nil {
		log.Println("数据格式转换失败")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误，评论创建失败"})
		return
	}

	dc := models.DocumentComment{
		DocumentID:   docId,
		ParentID:     nil,
		RootID:       nil,
		UserID:       contextData.UserId,
		Content:      contextData.CommentContent,
		Status:       "已发布",
		LikeCount:    0,
		DislikeCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsDeleted:    false,
		CreatedAtBy:  "",
		EditedAtBy:   "",
	}
	err = cc.commentDao.CreateComment(&dc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "系统错误，评论创建失败"})
		return
	}

	err = cc.commentDao.InsertCommentToRedis(dc)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"error": "系统错误，评论创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// 拉取顶级评论
func (cc *CommentController) GetDocumentRootComment(c *gin.Context) {
	docIdStr := c.Param("doc_id")
	pageStr := c.DefaultQuery("page", "0")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInsufficientStorage, gin.H{"error": "系统错误，评论信息拉取失败"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInsufficientStorage, gin.H{"error": "系统错误，评论信息拉取失败"})
		return
	}

	docId, err := strconv.ParseInt(docIdStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInsufficientStorage, gin.H{"error": "系统错误，评论信息拉取失败"})
		return
	}
	commentList, total, err := cc.commentDao.GetRootCommentsByDocumentID(docId, page, pageSize)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInsufficientStorage, gin.H{"error": "系统错误，评论信息拉取失败"})
		return
	}
	log.Println("评论数：%d", total)
	log.Println(commentList)
	var resultData []map[string]interface{}
	for _, comment := range commentList {
		have_children_comment, err := cc.commentDao.HasRepliesByCommentID(comment.ID)
		if err != nil {
			have_children_comment = false
		}
		tmp := map[string]interface{}{
			"comment_id":            strconv.FormatInt(comment.ID, 10),
			"comment_content":       comment.Content,
			"doc_id":                strconv.FormatInt(comment.DocumentID, 10),
			"user_id":               strconv.FormatInt(comment.UserID, 10),
			"nickname":              comment.User.Nickname,
			"last_updated_at":       comment.UpdatedAt,
			"comment_like_count":    comment.LikeCount,
			"have_children_comment": have_children_comment,
		}
		resultData = append(resultData, tmp)
	}
	c.JSON(http.StatusOK, gin.H{"comment_list": resultData})

}

// 根据顶级评论id获取子评论
