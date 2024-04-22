package service

import (
	"context"
	"dexbot/model"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var (
	db  *gorm.DB
	rdb *redis.Client
)

func GetDB() *gorm.DB {
	if db == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			GetConfig().Datasource.UserName,
			GetConfig().Datasource.Password,
			GetConfig().Datasource.Host,
			GetConfig().Datasource.Database)

		db_tmp, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			//Logger: gormlogger.Default.LogMode(gormlogger.Info),
		})
		if err != nil {
			panic("failed to connect database")
		}
		log.Println("数据库初始化成功...")
		db = db_tmp
		createTable(db_tmp)
	}
	return db
}

func createTable(database *gorm.DB) {
	//if err := database.AutoMigrate(&model.User{}); err != nil {
	//	log.Println("建表时出现错误", err)
	//}
	//if err := database.AutoMigrate(&model.Trade{}); err != nil {
	//	log.Println("建表时出现错误", err)
	//}
	if err := database.AutoMigrate(&model.WalletHistory{}); err != nil {
		log.Println("建表时出现错误", err)
	}
	if err := database.AutoMigrate(&model.GroupInfo{}); err != nil {
		log.Println("建表时出现错误", err)
	}
	if err := database.AutoMigrate(&model.ScheduleDelete{}); err != nil {
		log.Println("建表时出现错误", err)
	}
	//if err := database.AutoMigrate(&model.RedPacket{}); err != nil {
	//	log.Println("建表时出现错误", err)
	//}

	log.Println("建表成功...")
}

func GetRedis() *redis.Client {
	if rdb == nil {
		opts, err := redis.ParseURL("redis://localhost:6379/0")
		if err != nil {
			panic(err)
		}

		rdb = redis.NewClient(opts)
		if err = rdb.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
		log.Println("Redis 连接成功")
	}
	return rdb
}
