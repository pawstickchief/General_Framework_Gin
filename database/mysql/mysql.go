// database/mysql/mysql.go
package mysql

import (
	"fmt"
	"log"
	"sync"

	"General_Framework_Gin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var once sync.Once

// Init 初始化 MySQL 数据库连接
func Init() {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.AppConfig.Database.MySQL.User,
			config.AppConfig.Database.MySQL.Password,
			config.AppConfig.Database.MySQL.Host,
			config.AppConfig.Database.MySQL.Port,
			config.AppConfig.Database.MySQL.DBName,
		)

		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("无法连接到 MySQL 数据库: %v", err)
		}
	})
}

// Close 关闭 MySQL 数据库连接
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			log.Println("正在关闭数据库连接...")
			sqlDB.Close()
			log.Println("数据库连接已关闭")
		} else {
			log.Printf("数据库关闭失败: %v", err)
		}
	}
}
