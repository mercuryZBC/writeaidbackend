package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"yuqueppbackend/controllers"
	"yuqueppbackend/dao"
	"yuqueppbackend/db"
	"yuqueppbackend/util"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 设置 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许的跨域来源（可以是 *，但不推荐用于生产环境）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Content-Type", "Authorization"},           // 允许的请求头
		AllowCredentials: true,                                                // 是否允许携带凭证（如 Cookies）
	}))

	// 初始化 DAO 和 Controller
	docDao := dao.NewDocDao(db.GetDB())
	docController := controllers.NewDocumentController(docDao)
	dcDao := dao.NewCommentDAO(db.GetDB())
	dcController := controllers.NewCommentController(dcDao)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("register", controllers.Register)
		authGroup.POST("login", controllers.Login)
	}

	userGroup := r.Group("/api/user")
	userGroup.Use(util.AuthMiddleware())
	{
		userGroup.GET("getUserInfo", controllers.GetUserInfo)
		userGroup.POST("logout", controllers.Logout)
	}

	utilGroup := r.Group("/api/util")
	{
		utilGroup.GET("getCaptcha", controllers.GetCaptcha)

	}

	knowledgeGroup := r.Group("/api/knowledge")
	knowledgeGroup.Use(util.AuthMiddleware()) // 使用认知中间件
	{
		knowledgeGroup.POST("/createKnowledgeBase", controllers.CreateKnowledgeBase)
		knowledgeGroup.GET("/getKnowledgeBaseList", controllers.GetKnowledgeBaseList)
		knowledgeGroup.GET("/:kb_id", controllers.GetKnowledgeBaseDetail)
		knowledgeGroup.POST("/updateKnowledgeBase", controllers.UpdateKnowledgeBase)
		knowledgeGroup.POST("/deleteKnowledgeBase", controllers.DeleteKnowledgeBase)
	}

	documentGroup := r.Group("/api/document")
	documentGroup.Use(util.AuthMiddleware())
	{
		// 文档相关路由
		documentGroup.POST("/createDocument", docController.CreateDocumentHandler)
		documentGroup.GET("/getDocument/:doc_id", docController.GetDocumentByIDHandler)
		documentGroup.GET("/getDocumentListByKbId/:kb_id", docController.GetDocumentsByKnowledgeBaseIDHandler)
		documentGroup.PUT("/updateDocument/:doc_id", docController.UpdateDocumentHandler)
		documentGroup.DELETE("/deleteDocument/:doc_id", docController.DeleteDocumentByIDHandler)
		documentGroup.POST("/documents/:doc_id/view", docController.IncrementViewCountHandler)
		documentGroup.GET("/recentViewDocument", docController.GetRecentViewDocumentsHandler)
		documentGroup.GET("/recentEditDocument", docController.GetRecentEditDocumentsHandler)
		documentGroup.GET("/recentCommentDocument", docController.GetRecentCommentDocumentsHandler)
	}
	documentCommentGroup := r.Group("/api/comment")
	documentCommentGroup.Use(util.AuthMiddleware())
	{
		documentCommentGroup.POST("/createDocumentComment", dcController.CreateDocumentComment)
		documentCommentGroup.POST("/replyDocumentComment", dcController.ReplyDocumentComment)
		documentCommentGroup.GET("/getDocumentRootComment/:doc_id", dcController.GetDocumentRootComment)
	}

	return r
}
