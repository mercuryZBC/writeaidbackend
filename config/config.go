package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
	"os/user"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return err
	}
	return nil
}

func GetServerPort() string {
	return viper.GetString("server.port")
}

func GetDatabaseConfig() map[string]interface{} {
	return viper.GetStringMap("database")
}

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() error {
	// 数据库连接字符串
	dbConfig := GetDatabaseConfig()

	host := dbConfig["host"].(string)
	port := dbConfig["port"].(string)
	admin := dbConfig["user"].(string)
	password := dbConfig["password"].(string)
	dbname := dbConfig["dbname"].(string)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		admin, password, host, port, dbname)
	var err error

	// 连接数据库
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return err
	}

	// 自动迁移：自动创建/更新数据库表结构
	if err := DB.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database connected successfully")
	return nil
}
