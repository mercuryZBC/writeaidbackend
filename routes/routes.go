package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"yuqueppbackend/controllers"
	"yuqueppbackend/util"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 设置 CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                   // 允许的跨域来源（可以是 *，但不推荐用于生产环境）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的 HTTP 方法
		AllowHeaders:     []string{"Content-Type", "Authorization"},           // 允许的请求头
		AllowCredentials: true,                                                // 是否允许携带凭证（如 Cookies）
	}))

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
		knowledgeGroup.POST("/getKnowledgeBaseDetail", controllers.GetKnowledgeBaseDetail)
		knowledgeGroup.POST("/updateKnowledgeBase", controllers.UpdateKnowledgeBase)
		knowledgeGroup.POST("/deleteKnowledgeBase", controllers.DeleteKnowledgeBase)
	}

	return r
}
