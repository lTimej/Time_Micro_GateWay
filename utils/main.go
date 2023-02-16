package main

import (
	"context"
	"flag"
)

const (
	KEY_RBAC_REFRESH       = "rbac_handler_refresh" //强制刷新rbac的标识key
	KEY_RBAC_REFRESH_TIMES = "rbac_handler_times"   //累计cron执行检测次数key
	KEY_RBAC_ROLE_PREFIX   = "rbac_role:"           //角色对应权限的key前缀
	KEY_RBAC_METHOD        = "rbac_method"          //服务方法hash结构的key值
	KEY_RBAC_METHOD_PREFIX = "rbac_method:"         //服务方法hash结构的field前缀

	DB_USER = "xdrg"
	DB_PWD  = "xdrg123456"
	DB_ADDR = "127.0.0.1"
	DB_PORT = 3306
	DB_NAME = "emicro"

	REDIS_ADDR = "127.0.0.1"
	REDIS_PORT = 6379
)

var (
	InitTable string
	InitDb    string
)

func main() {
	flag.StringVar(&InitTable, "c", "", "init mysql table")
	flag.StringVar(&InitDb, "i", "", "init mysql database")
	flag.Parse()
	NewInit()
	if InitTable == "init_table" {
		DBMigrate()
		if InitDb == "init_db" {
			DBInit()
		}
	}
	if RED.Get(context.Background(), KEY_RBAC_REFRESH).Val() == "1" {
		err := RefreshRbacHandler()
		if err == nil {
			RED.Del(context.Background(), KEY_RBAC_REFRESH)
		}
		return
	}
	if RED.HLen(context.Background(), KEY_RBAC_METHOD).Val() <= 0 {
		RefreshRbacMethodsWithDb()
	}
	if RED.SCard(context.Background(), KEY_RBAC_ROLE_PREFIX+"1").Val() <= 0 {
		RefreshRbacRoleMethodsWithDb()
	}
	if RED.Get(context.Background(), KEY_RBAC_REFRESH_TIMES).Val() == "30" {
		if RefreshRbacHandler() != nil {
			return
		}
		RED.Set(context.Background(), KEY_RBAC_REFRESH_TIMES, 1, 0)
	} else {
		RED.Incr(context.Background(), KEY_RBAC_REFRESH_TIMES)
	}
}
