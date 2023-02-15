package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"liujun/Time_Micro_GateWay/server/common"
	"log"
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
func DBInit() {
	viper.SetConfigName("privileges") // name of config file (without extension)
	viper.SetConfigType("yml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("conf")       // path to look for the config file in 	// optionally look for config in the working directory
	err := viper.ReadInConfig()       // Find and read the config file
	if err != nil {                   // Handle errors reading the config file
		log.Printf("配置文件错误: %v", err)
		return
	}
	m := viper.AllSettings()
	tx := DB.Begin()
	for k, v := range m {
		p_method := Method{Name: k}
		if err := tx.Create(&p_method).Error; err != nil {
			log.Println("插入表method错误,err:", err)
			tx.Rollback()
			return
		}
		for _, handle := range v.([]interface{}) {
			for kk, vv := range handle.(map[string]interface{}) {
				c_method := Method{Name: kk, ParentId: p_method.Id}
				if err := tx.Create(&c_method).Error; err != nil {
					log.Println("插入表method错误,err:", err)
					tx.Rollback()
					return
				}
				for _, vvv := range vv.([]interface{}) {
					role := Role{}
					DB.Where("name = ?", vvv.(string)).First((&role))
					if role.Name == "" {
						role = Role{Name: vvv.(string)}
						if err := tx.Create(&role).Error; err != nil {
							log.Println("插入表role错误,err:", err)
							tx.Rollback()
							return
						}
					}
					role_method := RoleMethod{MethodId: c_method.Id, RoleId: role.RoleId}
					if err := tx.Create(&role_method).Error; err != nil {
						log.Println("插入表role_method错误,err:", err)
						tx.Rollback()
						return
					}
				}
			}
		}
	}
	tx.Commit()
	return
}
