package models

import (
	"fmt"
	"liujun/Time_Micro_GateWay/server/common"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB  *gorm.DB
	err error
)

func init() {
	db_host := common.Config.String("db_host")
	db_port, _ := common.Config.Int("db_port")
	db_name := common.Config.String("db_name")
	db_username := common.Config.String("db_username")
	db_password := common.Config.String("db_password")
	//root:liujun@tcp(127.0.0.1:3308)/micro_gateway?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, db_host, db_port, db_name)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //禁用表名复数
		},
	})
	if err != nil {
		log.Println("mysql init failed... ...")
		return
	}
	log.Println("mysql init success... ...")
}
