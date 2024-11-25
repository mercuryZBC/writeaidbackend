package main

import (
	"log"
	"yuqueppbackend/config"
	"yuqueppbackend/db"
	"yuqueppbackend/routes"
	"yuqueppbackend/util"
)

func main() {
	// 初始化配置

	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 初始化数据库单例对象
	db.GetDB()
	util.GetRedisClient()
	r := routes.SetupRouter()
	r.Run(config.GetServerPort())
}
