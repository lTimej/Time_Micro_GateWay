package main

import (
	"context"
	"errors"
	"log"
	"strconv"
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

type Method struct {
	Id       int    `gorm:"type:int(11) unsigned; primary_key; COMMENT '方法id'"`
	Name     string `gorm:"type:varchar(50); NOT NULL; COMMENT '方法名'"`
	ParentId int    `gorm:"type:int(11) unsigned; NOT NULL; DEFAULT: 0; COMMENT '所属微服务模块id'"`
}

type RoleMethod struct {
	MethodId int `gorm:"type:int(11) unsigned; primary_key; COMMENT '方法id'"`
	RoleId   int `gorm:"type:int(11) unsigned; primary_key; COMMENT '角色id'"`
}

func main() {
	NewInit()
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

func RefreshRbacMethodsWithDb() error {
	methods := []Method{}
	DB.Find(&methods)
	if len(methods) <= 0 {
		return errors.New("methods为空")
	}
	m := make(map[string]interface{})
	for _, method := range methods {
		m[KEY_RBAC_METHOD_PREFIX+method.Name] = method.Id
	}
	_, err := RED.HMSet(context.Background(), KEY_RBAC_METHOD, m).Result()
	return err
}

func RefreshRbacRoleMethodsWithDb() error {
	role_methods := []RoleMethod{}
	DB.Find(&role_methods)
	if len(role_methods) <= 0 {
		return errors.New("role_methods为空")
	}
	m := make(map[int][]string)
	for _, role_method := range role_methods {
		if v, ok := m[role_method.RoleId]; ok {
			m[role_method.RoleId] = append(v, strconv.Itoa(role_method.MethodId))
		} else {
			m[role_method.RoleId] = []string{strconv.Itoa(role_method.MethodId)}
		}
	}
	for role_id, method := range m {
		_, err := RED.SAdd(context.Background(), KEY_RBAC_ROLE_PREFIX+strconv.Itoa(role_id), method).Result()
		if err != nil {
			return errors.New("rbac_role_ Err:" + err.Error())
		}
	}
	return nil
}

func RefreshRbacHandler() error {
	err := RefreshRbacMethodsWithDb()
	if err != nil {
		log.Println("method_err:", err)
		return err
	}
	err = RefreshRbacRoleMethodsWithDb()
	if err != nil {
		log.Println("role_method_err:", err)
		return err
	}
	return nil
}
