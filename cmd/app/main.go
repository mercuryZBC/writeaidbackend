package main

import (
	"github.com/gin-gonic/gin"
	"yuqueppbackend/config"
	"yuqueppbackend/controllers"
	"yuqueppbackend/db"
)

func main() {

	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	// 初始化数据库单例对象
	db.GetDB()
	r := gin.Default()

	userGroup := r.Group("/api/user")
	{
		userGroup.GET("login", controllers.LoginPage)
		userGroup.POST("register", controllers.Register)
		userGroup.POST("login", controllers.Login)
	}
	r.Run(":8080")
}
