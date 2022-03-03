package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"sync"
	"time"
)

var DB *gorm.DB
var once sync.Once

// 初始化db
func InitDB(mysqlConfs map[string]interface{}) {
	once.Do(func() {
		fmt.Println("db init ....")

		url, ok := mysqlConfs["url"]
		if !ok {
			log.Fatal("mysql.url not find")
		}
		username, ok := mysqlConfs["username"]
		if !ok {
			log.Fatal("mysql.username not find")
		}
		password, ok := mysqlConfs["password"]
		if !ok {
			log.Fatal("mysql.password not find")
		}
		//dsn := "tao:taodyd123@tcp(127.0.0.1:3306)/ht-admin?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf("%s:%s@%s", username, password, url)
		fmt.Println(dsn)
		var err error = nil
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix: "gormv2_",
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("failed to connect database")
		}

		sqlDB, _ := DB.DB()
		// 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)
		// 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)
	})

}
