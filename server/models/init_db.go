package models

import (
	"context"
	"fmt"
	"liujun/Time_Micro_GateWay/server/common"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB  *gorm.DB
	RED *redis.Client
	err error
)

func init() {
	db_host := common.Config.String("db_host")
	db_port, _ := common.Config.Int("db_port")
	db_name := common.Config.String("db_name")
	db_username := common.Config.String("db_username")
	db_password := common.Config.String("db_password")
	//root:liujun@tcp(127.0.0.1:3308)/micro_gateway?charset=utf8mb4&parseTime=True&loc=Local
	// dsn := "root:liujun@tcp(127.0.0.1:3308)/micro_gateway?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, db_host, db_port, db_name)
	fmt.Println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //禁用表名复数
		},
	})
	if err != nil {
		log.Println("mysql init failed... ...err:", err)
		return
	}
	log.Println("mysql init success... ...")
}

func init() {
	config := common.Config
	redis_port, _ := config.Int("redis_port")
	redis_database, _ := config.Int("redis_database")
	redis_poolsize, _ := config.Int("redis_poolsize")
	redis_minidlecon, _ := config.Int("redis_minidlecon")
	RED = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.String("redis_host"), redis_port),
		Password:     config.String("redis_password"),
		DB:           redis_database,
		PoolSize:     redis_poolsize,
		MinIdleConns: redis_minidlecon,
	})
	ping, err := RED.Ping(context.Background()).Result()
	if err != nil {
		log.Println("init redis failed... ...", err)
		return
	}
	log.Println("redis init success... ...", ping)
}

func DBMigrate() {
	if err := DB.AutoMigrate(&User{}, &UserRole{}, &Role{}, &Method{}, &RoleMethod{}, &Service{}); err != nil {
		log.Println("数据迁移失败err:", err)
	} else {
		log.Println("数据迁移成功")
	}

}
