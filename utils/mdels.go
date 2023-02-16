package main

import "time"

// 用户表
type User struct {
	ID         int       `gorm:"type:bigint; NOT NULL"`
	Username   string    `valid:"Required;MinSize(4);MaxSize(12)"`
	Phone      string    `valid:"Mobile;Required"`
	Email      string    `valid:"Email;MaxSize(50)"`
	Password   string    `valid:"Required;MinSize(4);MaxSize(6)"`
	RePassword string    `gorm:"-" valid:"Required;MinSize(4);MaxSize(6)"`
	RegTime    time.Time `gorm:"type:datetime; DEFAULT: CURRENT_TIMESTAMP"`
}

// 用户角色表
type UserRole struct {
	RoleId int `gorm:"type:int(11) unsigned; NOT NULL; DEFAULT: 0; COMMENT '用户id'"`
	UserId int `gorm:"type:int(11) unsigned; primary_key; COMMENT '用户角色id'"`
}

// 方法表
type Method struct {
	Id       int    `gorm:"type:int(11) unsigned; primary_key; COMMENT '方法id'"`
	Name     string `gorm:"type:varchar(50); NOT NULL; COMMENT '方法名'"`
	ParentId int    `gorm:"type:int(11) unsigned; NOT NULL; DEFAULT: 0; COMMENT '所属微服务模块id'"`
}

// 角色表
type Role struct {
	RoleId int    `gorm:"type:int(11) unsigned; primary_key; COMMENT '角色id'"`
	Name   string `gorm:"type:varchar(50); NOT NULL ;COMMENT '角色名'"`
}

// 角色方法表
type RoleMethod struct {
	MethodId int `gorm:"type:int(11) unsigned; primary_key; COMMENT '方法id'"`
	RoleId   int `gorm:"type:int(11) unsigned; primary_key; COMMENT '角色id'"`
}

// grpc服务模块表
type Service struct {
	Id   int    `gorm:"type:int(11) unsigned; primary_key; COMMENT '服务模块id'"`
	Name string `gorm:"type:varchar(50); NOT NULL; COMMENT '服务模块名'"`
}
